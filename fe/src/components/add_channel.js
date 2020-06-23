const m = require("mithril");
const { getSingleError } = require("../common/error");
const Refresh = require("../common/Refresh");
const getLoadingView = require("./loading");
const UserChannels = require("./user_channels");
const { checkAuthAndExtract } = require("../actions/request");
const User = require("../actions/user");

const Config = require("../config");

const AddChannel = {
  current: {},
  state: "ready",
  setError: function(error) {
    this.current.error = error;
  },
  submit: function() {
    AddChannel.state = "in_progress";

      const channelUrlInput = document.getElementById("form-channel-url");
      const categoryInput = document.getElementById("form-channel-category");
      const categoryId = categoryInput[categoryInput.selectedIndex].id;

    return m
      .request({
        method: "POST",
        url: Config.api_url + "/channels",
        data: {
            "channel_url": channelUrlInput.value,
            "category_id": categoryId,
        },
        withCredentials: true
      })
      .then(result => {
        this.setError("");
        AddChannel.state = "ready";
        Refresh.posts = true;
        Refresh.userChannels = true;

        m.route.set("/");
      })
      .catch(e => {
        AddChannel.state = "ready";
        this.setError(getSingleError(e.message));
      });
  }
};

const AddNewChannelComponent = {
    user_categories: [],
    oninit: node => {
        return m.request({
            method: "GET",
            url: Config.api_url + "/users/current/categories",
            withCredentials: true,
            responseType: "json",
            extract: checkAuthAndExtract
        }).then(function(result) {
            AddNewChannelComponent.user_categories = result['response']['user_categories'];
        });
    },
  view: function(node) {
    const getLoading = () => {
      return (
          <div>
          {getLoadingView()}
          </div>
      );
    };

    let options = [];

      for (let c in AddNewChannelComponent.user_categories) {
          const category = AddNewChannelComponent.user_categories[c];
          options.push(<option id={category['id']}>{category['title']}</option>);
      }

    const getForm = () => {

      return (
          <div class="form-group">
            <div class="new-channel-error">{AddChannel.current.error}</div>
            <div class="col-lg-5">
              <input id="form-channel-url" type="text" placeholder="RSS Channel URL" autofocus class="form-control"></input>
              <div class="mt-2">
                  <select id="form-channel-category">
                    {options}
                  </select>
              </div>
              <button type="submit" class="btn btn-primary mt-2">Add</button>
            </div>
          </div>
      );
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
