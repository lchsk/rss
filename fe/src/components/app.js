const m = require("mithril");

const NavBar = require("./nav_bar");
const MainView = require("./main_view");
const User = require("../actions/user");
const Config = require("../config");
const LoginComponent = require("./login");
const Footer = require("./Footer");
const Refresh = require("../common/Refresh");
const Posts = require("../actions/posts");

const App = {
  oninit: node => {
    m.request({
      method: "GET",
      url: Config.api_url + "/users/current",
      withCredentials: true,
      responseType: "json",
      extract: function(xhr, options) {
        if (xhr.status === 401) {
          LoginComponent.signedIn = false;
          return false;
        } else if (xhr.status === 200) {
          const response = options.deserialize(xhr.response);
          LoginComponent.username = response['username'];
          console.log("here", options.deserialize(xhr.response));
          LoginComponent.signedIn = true;
          return true;
        }

        LoginComponent.signedIn = false;
        return false;
      }
    }).then(function(signedIn) {
    });
  },
  view: node => {
    if (Refresh.posts) {
      Refresh.posts = false;
      Refresh.postsIntervalId = setInterval(function(){
        const page = m.route.param("page");
        const id = m.route.param("id");
        const ret = Posts.loadChannels(node.attrs.type, page, id);
        ret.then(result => {
          clearInterval(Refresh.postsIntervalId);
          Refresh.postsIntervalId = null;
        })
      }, 1000);
    }

    if (Refresh.userChannels) {
      Refresh.userChannels = false;
      Refresh.userChannelsId = setInterval(function(){
        const ret = User.loadChannels();
        ret.then(result => {
          clearInterval(Refresh.userChannelsId);
          Refresh.userChannelsId = null;
        })
      }, 1000);
    }

    return (
        <div class="app">
          <NavBar>

          </NavBar>
          <MainView>
            {node.children}
          </MainView>
          <Footer></Footer>
        </div>
    );
  }
};

module.exports = App;
