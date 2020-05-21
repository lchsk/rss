var m = require("mithril");

const Posts = require("../actions/posts");
const getLink = require("./link");
const User = require("../actions/user");

const PostsList = {
  oninit: function(node) {
    const page = m.route.param("page");
    const id = m.route.param("id");

    Posts.loadChannels(node.attrs.type, page, id);
  },
  onupdate: function(node) {
    const page = m.route.param("page");
    const id = m.route.param("id");

    if (Posts.isRefreshNeeded(node.attrs.type, page, id)) {
      Posts.loadChannels(node.attrs.type, page, id);
    }
  },
  view: function(node) {
    let rows = [];

    const pagination = Posts.pagination;

    let prevButtonCls = ""
    let nextButtonCls = ""

    if (pagination.next === -1) {
      nextButtonCls = ".disabled";
    }

    if (pagination.prev === -1) {
      prevButtonCls = ".disabled";
    }

    let base_url = '';
    let title = '';

    if (node.attrs.type === "inbox") {
      base_url = "/";
      title = 'Inbox';
    } else if (node.attrs.type === "channel") {
      const id = m.route.param("id");
      base_url = "/channels/" + id;
    } else if (node.attrs.type === "category") {
      const id = m.route.param("id");
      base_url = "/categories/" + id;
      if (User.channelsByCategory[id] !== undefined) {
        title = User.channelsByCategory[id].categoryTitle;
      }
    }

    let prev = getLink(".btn .btn-primary" + prevButtonCls, base_url + "/?page=" + pagination.prev, "Prev")
    let next = getLink(".btn .btn-primary" + nextButtonCls, base_url + "/?page=" + pagination.next, "Next");

    for (let i = 0; i < Posts.data.length; i++) {
      const post = Posts.data[i];

      rows.push(
          <tr>
          <td>{post.title}</td>
          </tr>
      )
    }

    return (
        <div class="container-fluid">
        <h1>{title}</h1>
        <div class="row">
        {prev}{next}
        </div>
        <div class="row">
        <table>
        {rows}
      </table>
        </div>
      </div>
    )
  }
};

module.exports = PostsList;
