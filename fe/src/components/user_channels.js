var m = require("mithril");

const getLink = require("./link");
const defDict = require("../common/data_structures");
const User = require("../actions/user");

const UserChannels = {
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
      channelsHtml.push(
        getLink(
          ".list-group-item .list-group-item-action .active",
          "/categories/" + categoryId,
          categories.dict[categoryId]["categoryTitle"]
        )
      );

      for (const channel of c1) {
        channelsHtml.push(
          getLink(
            ".list-group-item .list-group-item-action",
            "/channels/" + channel["channel_id"],
            channel["channel_url"]
          )
        );
      }
    }

    return m("div", [m("div.list-group", channelsHtml)]);
  }
};

module.exports = UserChannels;
