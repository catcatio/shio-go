const xid = require('xid-js')

const newToken = () => {
    return xid.next()
}

const systemEvent = () => ({
    "replyToken": '00000000000000000000000000000000',
    "type": "message",
    "timestamp": Date.now(),
    "source": {
        "type": "user",
        "userId": "{{LineSenderID}}"
    },
    "message": {
        "id": "325708",
        "type": "text",
        "text": "Hello, world"
    }
})

const followEvent = (userId, replyToken = undefined) => ({
    "replyToken": replyToken || newToken(),
    "type": "follow",
    "timestamp": Date.now(),
    "source": {
        "type": "user",
        "userId": userId || "{{LineSenderID}}"
    }
})

const textMessageEvent = (message, userId, replyToken = undefined) => ({
    "replyToken": replyToken || newToken(),
    "type": "message",
    "timestamp": Date.now(),
    "source": {
        "type": "user",
        "userId": userId || "{{LineSenderID}}"
    },
    "message": {
        "id": `${Date.now()}`,
        "type": "text",
        "text": message
    }
})

module.exports = {
    systemEvent,
    followEvent,
    textMessageEvent
}
