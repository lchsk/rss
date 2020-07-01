import m from "mithril";

const SignUpComponent = require("./components/signup");
const Logout = require("./actions/logout");
const App = require("./components/app");
const AppNoAuth = require("./components/app_no_auth");
const MainViewWithChannelsList = require("./components/main_view_channels_list");
const PostsList = require("./components/posts_list");
const LoginComponent = require("./components/login");
const PostView = require("./components/post_view");
const AddNewChannelComponent = require("./components/add_channel");
const Account = require("./components/Account");

m.route.prefix("#!");

m.route(document.body, "/", {
  "/": {
    render: function() {
      m.route.set("/index");
      // return <App redirectTo="landing"></App>;
    }
  },
  "/index": {
    render: function() {
      return (
        <App>
          <MainViewWithChannelsList>
            <PostsList type="inbox"></PostsList>
          </MainViewWithChannelsList>
        </App>
      );
    }
  },
  "/channels/new": {
    render: function() {
      return (
        <App>
          <MainViewWithChannelsList>
            <AddNewChannelComponent></AddNewChannelComponent>
          </MainViewWithChannelsList>
        </App>
      );
    }
  },
  "/channels/:id": {
    render: function() {
      return (
        <App>
          <MainViewWithChannelsList>
            <PostsList type="channel"></PostsList>
          </MainViewWithChannelsList>
        </App>
      );
    }
  },
  "/categories/:id": {
    render: function() {
      return (
        <App>
          <MainViewWithChannelsList>
            <PostsList type="category"></PostsList>
          </MainViewWithChannelsList>
        </App>
      );
    }
  },
  "/posts/:id": {
    render: function() {
      return (
        <App>
          <MainViewWithChannelsList>
            <PostView></PostView>
          </MainViewWithChannelsList>
        </App>
      );
    }
  },
  "/login": {
    render: function() {
      return (
        <AppNoAuth>
          <LoginComponent></LoginComponent>
        </AppNoAuth>
      );
    }
  },
  "/logout": {
    render: function() {
      return <Logout></Logout>;
    }
  },
  "/signup": {
    render: () => {
      return (
        <AppNoAuth>
          <SignUpComponent></SignUpComponent>
        </AppNoAuth>
      );
    }
  },
  "/account": {
    render: () => {
      return (
        <App>
          <Account></Account>
        </App>
      );
    }
  }
});
