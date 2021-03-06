const m = require("mithril");

const Config = require("../config");
const User = require("./user");

const { checkAuthAndExtract } = require("./request");

const PostsSource = {
  page: 1,
  type: "inbox",
  id: ""
};

const Posts = {
  data: {},
  pagination: {},
  source: PostsSource,

  isRefreshNeeded: (type, page, id) => {
    page = page || 1;
    id = id || "";

    return (
      page !== Posts.source.page ||
      type !== Posts.source.type ||
      id !== Posts.source.id
    );
  },

  loadChannels: (type, page, id) => {
    page = page || 1;
    id = id || "";

    Posts.source.page = page;
    Posts.source.type = type;
    Posts.source.id = id;

    let url;

    if (type === "inbox") {
      url = Config.api_url + "/posts-inbox?page=" + page;
    } else if (type === "channel") {
      url =
        Config.api_url + "/posts/channels/" + Posts.source.id + "?page=" + page;
    } else if (type === "category") {
      url =
        Config.api_url +
        "/posts/categories/" +
        Posts.source.id +
        "?page=" +
        page;
    }

    return (
      m
        .request({
          method: "GET",
          url: url,
          withCredentials: true,
          responseType: "json",
          extract: checkAuthAndExtract
        })
        .then(result => {
          console.log("result", result);
          // User.defaultSuccess();

          const { response } = result;

          console.log(response);

          Posts.data = response.posts;
          Posts.pagination = response.pagination;
        })
        // , (error) => {
        // console.log("yo error", error, error.code);
        // })
        .catch(function(e) {
          Posts.data = null;
          Posts.pagination = null;
          // console.log("errrr", e, e.code, Object.keys(e));
        })
    );
  }
};

module.exports = Posts;
