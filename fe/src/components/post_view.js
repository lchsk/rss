var m = require("mithril");
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
      // extract: function(xhr, options) {
      //   if (xhr.status === 401) {
      //     LoginComponent.signedIn = false;
      //     return false;
      //   }

      //   LoginComponent.signedIn = true;
      //   return true;
      // }
    }).then(function(result) {
      console.log(result);
      PostView.postData = result.response;
      // console.log("signed-in", signedIn);
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
          <a href="{post.url}">{post.url}</a>
        </div>
        <div>{post.description}</div>
        <div>{post.content}</div>
      </div>
    );
  }
};

module.exports = PostView;
