#!/usr/bin/env node

"use strict"

const request = require('request');
const path = require('path');
const tar = require('tar');
const zlib = require('zlib');
const mkdirp = require('mkdirp');
const fs = require('fs');
const os = require('os');

// Mapping from Node's `process.arch` to Golang's `$GOARCH`
const ARCH_MAPPING = {
    "x32": "32bit",
    "x64": "64bit",
    "arm": "arm",
    "arm64": "arm64",
};

// Mapping between Node's `process.platform` to Golang's
const PLATFORM_MAPPING = {
    "darwin": "macOS",
    "linux": "linux",
    "win32": "windows",
};

function validateConfiguration(packageJson) {
    if (!packageJson.version) {
        return "'version' property must be specified";
    }

    if (!packageJson.goBinary || typeof (packageJson.goBinary) !== "object") {
        return "'goBinary' property must be defined and be an object";
    }

    if (!packageJson.goBinary.name) {
        return "'name' property is necessary";
    }

    if (!packageJson.goBinary.path) {
        return "'path' property is necessary";
    }

    if (!packageJson.goBinary.url) {
        return "'url' property is required";
    }
}

function parsePackageJson() {
    if (!(os.arch() in ARCH_MAPPING)) {
        console.error("Installation is not supported for this architecture: " + process.arch);
        return;
    }

    if (!(os.platform() in PLATFORM_MAPPING)) {
        console.error("Installation is not supported for this platform: " + process.platform);
        return
    }

    const packageJsonPath = path.join(".", "package.json");
    if (!fs.existsSync(packageJsonPath)) {
        console.error("Unable to find package.json. " +
            "Please run this script at root of the package you want to be installed");
        return
    }

    var packageJson = JSON.parse(fs.readFileSync(packageJsonPath));
    var error = validateConfiguration(packageJson);
    if (error && error.length > 0) {
        console.error("Invalid package.json: " + error);
        return
    }

    // We have validated the config. It exists in all its glory
    var binPath = packageJson.goBinary.path;
    var url = packageJson.goBinary.url;
    var wioName = packageJson.goBinary.wioName;
    var version = packageJson.version;
    if (version[0] === 'v') version = version.substr(1);  // strip the 'v' if necessary v0.0.1 => 0.0.1

    if (process.platform === "win32") {
        wioName = wioName.replace(/{{format}}/g, ".exe");
    } else {
        wioName = wioName.replace(/{{format}}/g, "");
    }

    // Interpolate variables in URL, if necessary
    url = url.replace(/{{arch}}/g, ARCH_MAPPING[os.arch()]);
    url = url.replace(/{{platform}}/g, PLATFORM_MAPPING[os.platform()]);
    url = url.replace(/{{version}}/g, version);

    if (process.platform === "win32") {
        url = url.replace(/{{format}}/g, "zip");
    } else {
        url = url.replace(/{{format}}/g, "tar.gz");
    }

    return {
        binPath: binPath,
        wioName: wioName,
        url: url,
        version: version
    }
}

const INVALID_INPUT = "Invalid inputs";
function install(callback) {
    var opts = parsePackageJson();
    if (!opts) return callback(INVALID_INPUT);

    var req = request({ uri: opts.url });

    if (process.platform === "win32") {
        console.log("Downloading and unzipping build files");

        mkdirp.sync(opts.binPath);

        var file = fs.createWriteStream(opts.binPath + "/wio.zip");
        req.on('error', callback.bind(null, "Error downloading from URL: " + opts.url));
        req.on('response', function (res) {
            if (res.statusCode !== 200) {
                return callback("Error downloading Wio binary. HTTP Status Code: " + res.statusCode);
            }

            req.on('data', function (chunk) {
                file.write(chunk);
            }).on('end', function () {
                file.end();
            });
        });

        req.on('end', function () {
            console.log('downloaded wio.zip!');

            const decompress = require('decompress');

            decompress(opts.binPath + '/wio.zip', opts.binPath).then(files => {
                console.log('unzipped wio.zip!');

                fs.rename(opts.binPath + "/" + opts.wioName, opts.binPath + '/wio', function(err) {
                    if ( err ) console.log('ERROR: ' + err);
                });
            });
        });

    } else {
        mkdirp.sync(opts.binPath);

        console.log("Downloading and extracting build files");

        var ungz = zlib.createGunzip();
        var untar = tar.Extract({ path: opts.binPath });

        untar.on('end', function () {
            console.log("Wio binary successfuly downloaded and extracted")
        });

        req.on('error', callback.bind(null, "Error downloading from URL: " + opts.url));
        req.on('response', function (res) {
            if (res.statusCode !== 200) return callback("Error downloading Wio binary. HTTP Status Code: " + res.statusCode);
            req.pipe(ungz).pipe(untar);
        });
    }
}

var myCallback = function (data) {
    console.log(data);
    process.exit(1)
};

// call the install
install(myCallback);
