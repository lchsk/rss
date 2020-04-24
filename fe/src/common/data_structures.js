


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

module.exports = defDict;
