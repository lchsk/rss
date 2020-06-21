const m = require("mithril");
const Config = require("../config");
const { checkAuthAndExtract } = require("../actions/request");

const PostView = {
  postData: {},
  oninit: node => {
    const postId = m.route.param("id");

    m.request({
      method: "GET",
      url: Config.api_url + "/posts/" + postId,
      withCredentials: true,
      responseType: "json",
      extract: checkAuthAndExtract
    }).then(function(result) {
      console.log(result);
      PostView.postData = result.response;

      // TODO: This request should be conditioned on the current status of the post (need to pull that status)
        m.request({
            method: "PATCH",
            url: Config.api_url + "/posts/" + postId,
            withCredentials: true,
            responseType: "json",
            data: {'status': 'read'},
            extract: checkAuthAndExtract
        }).then(function(result) {
            console.log(result);
        });
    });
  },
  view: () => {
    const post = PostView.postData;
    return (
      <div>
        <h1>{post.title}</h1>
        <div>
          {post.author_name} {post.author_email}
        </div>
        <div>{new Date(post.pub_at).toLocaleString()}</div>
        <div>
          <a href={post.url}>{post.url}</a>
        </div>
        <div>{post.description}</div>
        <div>{post.content}</div>
      </div>
    );
  }
};

module.exports = PostView;
