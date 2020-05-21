import m from "mithril";

const SignUpComponent = require("./components/signup");
const Logout = require("./actions/logout");
const App = require("./components/app");
const AppNoAuth = require("./components/app_no_auth");
const MainViewWithChannelsList = require("./components/main_view_channels_list");
const PostsList = require("./components/posts_list");
const User = require("./actions/user");
const Posts = require("./actions/posts");
const LoginComponent = require("./components/login");
const AddNewChannelComponent = require("./components/add_channel");
const { getErrorMessage, getSingleError } = require("./common/error");

User.loadChannels();

m.route.prefix("#!");

m.route(document.body, "/", {
  "/": {
    render: function() {
      return (
          <App>
          <MainViewWithChannelsList>
          <PostsList type="inbox">
          </PostsList>

        </MainViewWithChannelsList>
        </App>
      )
    }
  },
  "/channels/new": {
    render: function() {
      return m(App, m(MainViewWithChannelsList, m(AddNewChannelComponent)));
    }
  },
  "/channels/:id": {
	render: function() {
      return (
          <App>
          <MainViewWithChannelsList>
          <PostsList type="channel">
          </PostsList>

        </MainViewWithChannelsList>
        </App>
      )
	}
  },
  "/categories/:id": {
	render: function() {
      return (
          <App>
          <MainViewWithChannelsList>
          <PostsList type="category">
          </PostsList>

        </MainViewWithChannelsList>
        </App>
      )
	}
  },
  "/login": {
    render: function() {
      return m(AppNoAuth, m(LoginComponent));
    }
  },
  "/logout": {
    render: function() {
      return m(Logout);
    }
  },
  "/signup": {
    render: () => {
      return m(AppNoAuth, m(SignUpComponent));
    }
  }
});
