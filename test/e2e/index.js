const fs = require('fs');
const path = require('path');
const starman = require('@rungsikorn/starman').default
const commander = require('commander');

let postmanEnvironmentForFixtures = ({
    APIEndpoint: "http://localhost:30001",
    LineChatWebhook: "/chat/line",
    LineIgnoreSignature: true,
    LineSenderID: "U91eeaf62d901234567890123456789ab",
})

const start = () => {
    commander
        .option('-c, --collection <name>', 'specify comma separated collections to test', '')

    commander.on('--help', function(){
        console.log('')
        console.log('Examples:');
        console.log('  $ node runner');
        console.log('  $ node runner -c profile,receiver');
        console.log('  $ node runner --collection profile,receiver');
    });

    commander.parse(process.argv);

    let collections = commander.collection.split(",").join("|")

    const steps = []
    fs.readdirSync(path.join(__dirname, "./collections")).forEach(f => {
        if (RegExp(`(${collections})\.e2e`).test(f)) {
            steps.push(require(`./collections/${f}`))
        }
    })

    starman(steps, postmanEnvironmentForFixtures, {
        outputDir: path.resolve(__dirname, "../postman"),
        collectionName: "shio-go",
        environmentName: "shio-go",
    })
}

start()
