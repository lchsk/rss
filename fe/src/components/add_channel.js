var m = require("mithril");
var { getErrorMessage, getSingleError } = require("../common/error");
const getLoadingView = require("./loading");

const Config = require("../config");

var AddChannel = {
  current: {},
  state: "ready",
  setError: function(error) {
    this.current.error = error;
  },
  submit: function() {
    AddChannel.state = "in_progress";

    return m
      .request({
        method: "POST",
        url: Config.api_url + "/channels",
        data: AddChannel.current,
        withCredentials: true
      })
      .then(result => {
        this.setError("");
        AddChannel.state = "ready";
        m.route.set("/index");
      })
      .catch(e => {
        AddChannel.state = "ready";
        this.setError(getSingleError(e.message));
      });
  }
};

var AddNewChannelComponent = {
  view: function(node) {
    const getLoading = () => {
      return m("div", getLoadingView());
    };

    const getForm = () => {
      return m("div.form-group", [
        m("div#new-channel-error", AddChannel.current.error),
        m("div.col-lg-5", [
          m(
            "input[type=text][placeholder=RSS Channel URL] autofocus .form-control .form-control-lg",
            {
              autofocus: true,
              oninput: m.withAttr("value", value => {
                AddChannel.current.channel_url = value;
              })
            }
          ),
          m("small.form-text .text-muted", "Some helpful info"),
          m("div", { style: { paddingTop: "8px" } }),
          m("button[type=submit] .btn .btn-lg .btn-primary", "Add")
        ])
      ]);
    };

    const content =
      AddChannel.state === "in_progress" ? getLoading() : getForm();

    return m("div", [
      m("h1", "Add new channel"),
      m(
        "form",
        {
          onsubmit: e => {
            e.preventDefault();
            AddChannel.submit();
          }
        },
        [content]
      )
    ]);
  }
};

module.exports = AddNewChannelComponent;
