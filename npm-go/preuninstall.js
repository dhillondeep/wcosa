
const path = require('path')
const fs = require('fs')

// Mapping from Node's `process.arch` to Golang's `$GOARCH`
const ARCH_MAPPING = {
    "ia32": "32-bit",
    "x64": "64-bit",
    "arm": "arm"
};

// Mapping between Node's `process.platform` to Golang's 
const PLATFORM_MAPPING = {
    "darwin": "darwin",
    "linux": "linux",
    "win32": "windows",
    "freebsd": "freebsd"
};

function parsePackageJson() {
    const packageJsonPath = path.join(".", "package.json");
    if (!fs.existsSync(packageJsonPath)) {
        console.error("Unable to find package.json. " +
            "Please run this script at root of the package you want to be installed");
        return
    }

    let packageJson = JSON.parse(fs.readFileSync(packageJsonPath));

    // We have validated the config. It exists in all its glory
    let binName = packageJson.goBinary.name;
    let binPath = packageJson.goBinary.path;
    let url = packageJson.goBinary.url;
    let version = packageJson.version;
    if (version[0] === 'v') version = version.substr(1);  // strip the 'v' if necessary v0.0.1 => 0.0.1

    // Binary name on Windows has .exe suffix
    if (process.platform === "win32") {
        binName += ".exe"
    }

    // Interpolate variables in URL, if necessary
    url = url.replace(/{{arch}}/g, ARCH_MAPPING[process.arch]);
    url = url.replace(/{{platform}}/g, PLATFORM_MAPPING[process.platform]);
    url = url.replace(/{{version}}/g, version);

    if (process.platform === "win32") {
        url = url.replace(/{{format}}/g, "zip");
    } else {
        url = url.replace(/{{format}}/g, "tar.gz");
    }

    url = url.replace(/{{bin_name}}/g, binName);
    url = url.replace(/{{bin_name}}/g, binName);


    return {
        binName: binName,
        binPath: binPath,
        url: url,
        version: version
    }
}

function rename(callback) {
    let opts = parsePackageJson();
    if (!opts) return callback(INVALID_INPUT);

    if (process.platform === "win32") {
        fs.rename(opts.binPath + '/wio.exe', opts.binPath + '/wio', function(err) {
            if ( err ) console.log('ERROR: ' + err);
        });
    }
}

var myCallback = function (data) {
    console.log(data);
    process.exit(1)
};

rename(myCallback)
