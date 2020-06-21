const m = require("mithril");

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

    console.log("posts", Posts);
    const pagination = Posts.pagination;

    if (pagination === null) {
      return;
    }

    let prevButtonCls = "";
    let nextButtonCls = "";

    if (pagination.next === -1) {
      nextButtonCls = ".disabled";
    }

    if (pagination.prev === -1) {
      prevButtonCls = ".disabled";
    }

    let base_url = "";
    let title = "";

    if (node.attrs.type === "inbox") {
      base_url = "/";
      title = "Inbox";
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

    let prev = getLink(
      ".btn .btn-dark .together-left" + prevButtonCls,
      base_url + "?page=" + pagination.prev,
      "◄"
    );
    let next = getLink(
      ".btn .btn-dark .together-right" + nextButtonCls,
      base_url + "?page=" + pagination.next,
      "►"
    );

    for (let i = 0; i < Posts.data.length; i++) {
      const post = Posts.data[i];

      let postCls = "";

      if (post.status === "unread") {
        postCls = "unread-post";
      }
      rows.push(
        <tr class={postCls}>
          <td class="p-0">{getLink(".text-dark", "/posts/" + post.id, post.title)}</td>
          <td class="text-muted">{new Date(post.pub_at).toLocaleString()}</td>
        </tr>
      );
    }

    return (
      <div class="container-fluid">
        <div class="row">
          <h1>{title}</h1>
        </div>

        <div class="row mb-3">
          {prev}
          {next}
        </div>

        <div class="row">
          <table class="table table-borderless table-sm">{rows}</table>
        </div>
      </div>
    );
  }
};

module.exports = PostsList;
