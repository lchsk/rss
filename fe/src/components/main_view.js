const m = require("mithril");

const MainView = {
  view: node => {
    return (
        <div class="main-view">
          {node.children}
        </div>
    );
  }
};

module.exports = MainView;
