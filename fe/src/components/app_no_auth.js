var m = require("mithril");

const MainView = require("./main_view");

const AppNoAuth = {
  view: (node) => {
    return m(".app", [
      m(MainView, node.children),
    ]);
  }
};


module.exports = AppNoAuth;
