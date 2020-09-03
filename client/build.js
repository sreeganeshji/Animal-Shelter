const exec = require('child_process').exec;
const os = require('os');
const puts = (error, stdout, stderr) => console.log(stdout);

if (os.type() === 'Linux') {
    exec("cp -R build ../back-end", puts);
    exec("cp -R ../back-end $GOPATH/src", puts);
}
else if (os.type() === 'Darwin') {
    exec("cp -a build ../back-end", puts);
    exec("cp -a ../back-end $GOPATH/src")
}
else {
    // assume windows
    exec("xcopy /e /i build ../back-end", puts);
    exec("xcopy /e /i ../build $GOPATH/src", puts);
}