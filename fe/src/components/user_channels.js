const m = require("mithril");

const getLink = require("./link");
const defDict = require("../common/data_structures");
const User = require("../actions/user");

const UserChannels = {
  oninit: function(node) {
    console.log("user channels oninit");
    User.loadChannels();
  },
  onupdate: function(node) {
    // Detect when channels might have been changed and force refresh
    // User.loadChannels();
  },
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

    let channelsHtml = [];

    const link = (
        <p class="m-0">
          {getLink(
              ".text-dark",
              "/",
              "Inbox"
          )}
        </p>
    );
    channelsHtml.push(link);

    for (const [categoryId, c1] of Object.entries(channels.dict)) {
      const link = (
          <p class="m-0">
            {getLink(
                 ".text-dark",
               "/categories/" + categoryId,
               categories.dict[categoryId]["categoryTitle"]
            )}
          </p>
      );
      channelsHtml.push(link);

      for (const channel of c1) {
        const link = (
            <p class="m-0 ml-3">
              {getLink(
              ".text-dark",
              "/channels/" + channel["channel_id"],
              channel["channel_title"]
              )}
            </p>
        );
        channelsHtml.push(link);
      }
    }

    return (
        <div>
          <div class="text-left">
            {channelsHtml}
          </div>
        </div>
    );
  }
};

module.exports = UserChannels;
