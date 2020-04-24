var m = require("mithril");

const MainView = {
  view: node => {
    return m(".main-view", node.children);
  }
};

module.exports = MainView;
