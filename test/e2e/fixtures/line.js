const {followEvent, textMessageEvent} = require('./lineEvents')
const {StarmanRequestStep} = require('@rungsikorn/starman')
const expectation = require('./expect')

const sendEvents = (options, ...events) => {
    let reqBody = {"events": events}
    return new StarmanRequestStep(options.name || 'Send events')
        .AddHeader('Content-Type', 'application/json')
        .AddHeader('x-shio-debug', "{{LineIgnoreSignature}}")
        .Post(`{{APIEndpoint}}{{LineChatWebhook}}`)
        .AddBody(reqBody)
        .AddPreRequest(options.preRequest || expectation.Nothing)
        .AddTest(options.expect || expectation.Nothing)
}

const sendFollowEvent = (options, userId = undefined) => {
    options.name = options.name || "send follow event"
    return sendEvents(Object.assign(options), followEvent(userId))
}

const sendTextMessageEvent = (options, message) => {
    options.name = options.name || "send text message event"
    return sendEvents(Object.assign(options), textMessageEvent(message))
}

module.exports = {
    sendEvents,
    sendFollowEvent,
    sendTextMessageEvent,
}
