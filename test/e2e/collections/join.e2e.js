const lineFixtures = require('../fixtures/line')
const lineEvents = require('../fixtures/lineEvents')
const expect = require('../fixtures/expect')

const optionsExpectBadRequest = (options) => {
    options = options || {}
    options.expect = expect.ResponseStatus(400, "err: Bad Request")
    return options
}

const optionsExpectOK = (options) => {
    options = options || {}
    options.expect = expect.ResponseStatus(200, "OK")
    return options
}

module.exports = function (runner) {
    return runner("events", [
        lineFixtures.sendEvents(optionsExpectBadRequest({name: "send empty event"})),
        lineFixtures.sendEvents(optionsExpectOK({name: "send system event"}), lineEvents.systemEvent())
    ])
}
