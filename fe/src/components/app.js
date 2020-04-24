var m = require("mithril");

const NavBar = require("./nav_bar");
const MainView = require("./main_view");

const App = {
  view: (node) => {
    return m(".app", [
      m(NavBar),
      m(MainView, node.children),
    ]);
  }
};


module.exports = App;
