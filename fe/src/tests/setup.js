const o = require("ospec");
const jsdom = require("jsdom");

var dom = new jsdom.JSDOM("", {
    // So we can get `requestAnimationFrame`
    pretendToBeVisual: true,
})

global.window = dom.window
global.document = dom.window.document
global.requestAnimationFrame = dom.window.requestAnimationFrame

require("mithril")

o.after(function() {
    dom.window.close()
})
