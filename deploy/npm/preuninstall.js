
const path = require('path');
const fs = require('fs');
const os = require('os');

function parsePackageJson() {
    const packageJsonPath = path.join(".", "package.json");
    if (!fs.existsSync(packageJsonPath)) {
        console.error("Unable to find package.json. " +
            "Please run this script at root of the package you want to be installed");
        return
    }

    let packageJson = JSON.parse(fs.readFileSync(packageJsonPath));

    return {
        binName: packageJson.goBinary.name,
        binPath: packageJson.goBinary.path,
        url: "",
        version: ""
    }
}

function rename(callback) {
    var opts = parsePackageJson();
    if (!opts) return callback(INVALID_INPUT);

    if (os.platform() === "win32") {
        fs.rename(opts.binPath + '/wio.exe', opts.binPath + '/' + opts.binName, function(err) {
            if ( err ) console.log('ERROR: ' + err);
        });
    }
}

var myCallback = function (data) {
    console.log(data);
    process.exit(1)
};

rename(myCallback);

