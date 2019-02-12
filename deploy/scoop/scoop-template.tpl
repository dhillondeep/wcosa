{
    "version": "{{version}}",
    "architecture": {
        "32bit": {
            "url": "https://github.com/wio/wio/releases/download/v{{version}}/{{exec32bit}}.{{extension}}",
            "hash": "{{checksum32bit}}",
            "bin": [ ["{{exec32bit}}.exe", "wio"] ]
        },
        "64bit": {
            "url": "https://github.com/wio/wio/releases/download/v{{version}}/{{exec64bit}}.{{extension}}",
            "hash": "{{checksum64bit}}",
            "bin": [ ["{{exec64bit}}.exe", "wio"] ]
        }
    },
    "homepage": "https://github.com/wio/wio",
    "license": "MIT",
    "description": "An IOT Development Environment"
}
