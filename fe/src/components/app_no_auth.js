const m = require("mithril");

const MainView = require("./main_view");
const Footer = require("./Footer");

const AppNoAuth = {
  view: node => {
    return (
        <div class="app">
          <MainView>
              {node.children}
          </MainView>
            <Footer></Footer>
        </div>
    )
  }
};

module.exports = AppNoAuth;
