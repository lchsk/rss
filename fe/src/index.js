import m from "mithril";

const SignUpComponent = require("./signup");
const LoginComponent = require("./login");
const LogoutComponent = require("./logout");
const User = require("./user");
const AddNewChannelComponent = require("./add_channel");
var { getErrorMessage, getSingleError } = require("./error");

User.loadChannels();

const UserComponent = {
  oninit: node => {
    User.load();
  },
  view: node => {}
};

function defDict(type) {
  var dict = {};
  return {
    get: function(key) {
      if (!dict[key]) {
        dict[key] = type.constructor();
      }
      return dict[key];
    },
    dict: dict
  };
}

var UserChannels = {
  // oninit: User.loadChannels(),
  view: node => {
    var channels = defDict([]);
    var categories = defDict({});

    for (let i = 0; i < User.channels.length; i++) {
      const categoryId = User.channels[i]["category_id"];
      const categoryTitle = User.channels[i]["category_title"];

      if (categoryId) {
        categories.get(categoryId)["categoryTitle"] = categoryTitle;
      } else {
        categories.get(categoryId)["categoryTitle"] = "Without category";
      }

      channels.get(categoryId).push(User.channels[i]);
    }

    var channelsHtml = [];

    for (const [categoryId, c1] of Object.entries(channels.dict)) {
      console.log(categoryId);
      channelsHtml.push(
        getLink(
          ".list-group-item .list-group-item-action .active",
          "/user",
          categories.dict[categoryId]["categoryTitle"]
        )
      );

      for (const channel of c1) {
        console.log("\t" + channel["channel_url"]);
        channelsHtml.push(
          getLink(
            ".list-group-item .list-group-item-action",
            "/user",
            channel["channel_url"]
          )
        );
      }
    }

    return m("div", [m("div.list-group", channelsHtml)]);
  }
};

function getLink(classes, href, anchor) {
  return m(
    "a" + classes,
    {
      href: href,
      oncreate: m.route.link
    },
    anchor
  );
}

var Layout = {
  view: function(node) {
    return m(".container-fluid", [
      m(".row.top-bar", [
        m(".col-sm-3", [m("div", "hello")]),
        m(".col-sm-6", [m("div", "")]),
        m(".col-sm-3 .text-right", [m("div", "user")])
      ]),
      m(".row", [
        m(".col-sm .text-center", [
          getLink(".btn .btn-primary", "/channels/new", "Add new channel"),

          m("hr"),

          m(UserChannels)
          // m('div', [
          //   m('div.list-group', [
          //     getLink(".list-group-item .list-group-item-action", "/user", "Option 1"),
          //     getLink(".list-group-item .list-group-item-action", "/user", "Option 2"),
          //     getLink(".list-group-item .list-group-item-action", "/user", "Option 3"),
          //   ])
          // ]),
        ]),
        m(".col-sm-9", node.children)
      ])
    ]);
  }
};

var BoxLayout = {
  view: function(node) {
    return m(".text-center", [
      m(".center", [
        m(".top", "Logo"),
        m("", node.children),
        m(".bottom", "Forgot password? SignUP")
      ])
    ]);
  }
};

m.route.prefix("#!");

m.route(document.body, "/", {
  "/": {
    render: function() {
      console.log("route /");
      return m(Layout);
    }
  },
  "/channels/new": {
    render: function() {
      console.log("new channel");
      return m(Layout, m(AddNewChannelComponent));
    }
  },
  "/login": {
    render: function() {
      return m(BoxLayout, m(LoginComponent));
    }
  },
  "/logout": {
    render: function() {
      return m(LogoutComponent);
    }
  },
  "/signup": {
    render: () => {
      return m(SignUpComponent);
    }
  },
  "/user": {
    render: function() {
      return m(UserComponent);
    }
  }
});
