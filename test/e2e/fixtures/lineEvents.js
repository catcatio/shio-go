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
        "userId": "u206d25c2ea6bd87c17655609a1c37cb8"
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
        "userId": userId
    }
})

module.exports = {
    systemEvent,
    followEvent
}
