const webpack = require("webpack");

module.exports = {
  entry: "./src/index.js",
  module: {
    rules: [
      {
        test: /\.js$/,
        exclude: /\/node_modules\//,
        use: {
          loader: "babel-loader"
        }
      }
    ]
  },
  resolve: {
    extensions: ["*", ".js", ".jsx"]
  },
  output: {
    path: __dirname + "/dist/assets/",
    // publicPath: "/",
    filename: "bundle.js"
  },
  plugins: [new webpack.HotModuleReplacementPlugin()],
  devServer: {
    contentBase: "./dist",
    writeToDisk: filePath => {
      return /bundle.js$/.test(filePath);
    },
    // writeToDisk: true,
    hot: true
  }
};
