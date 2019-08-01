const template = (t) => {
    const func = new Function(`with(this) { return \`${t}\`; }`);
    return {
        fill: (d) => func.call(d)
    }
}

const makePostmanFunc = (func, data) => {
    return eval(template(func.toString()).fill(data))
}

const ResponseStatus = (statusCode, message = '', title = "should return") => {
    let a = (pm) => {
        pm.test('${title} ${statusCode}', function () {
            pm.response.to.have.status(parseInt('${statusCode}', 10));
            if ('${message}') {
                pm.expect(pm.response.text()).to.equal('${message}')
            }
        });
    }
    return makePostmanFunc(a, {statusCode, title, message})
}

const Nothing = () => {}

module.exports = {
    makePostmanFunc,
    ResponseStatus,
    Nothing
}
