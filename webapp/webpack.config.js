const path = require('path');
const PLUGIN_ID = require('../plugin.json').id;

module.exports = {
  entry: [
    './src/index.jsx',
  ],
  resolve: {
    modules: [
      'src',
      'node_modules',
      path.resolve(__dirname),
    ],
    extensions: ['*', '.js', '.jsx'],
  },
  module: {
    rules: [
      {
        test: /\.(js|jsx)$/,
        exclude: /node_modules/,
        use: {
          loader: 'babel-loader',
          options: {
            presets: ['@babel/preset-react',
              [
                "@babel/preset-env",
                {
                  "modules": "commonjs",
                  "targets": {
                    "node": "current"
                  }
                }
              ]
            ],
          },
        },
      },
      {
        test: /\.scss$/,
        use: [
          'style-loader',
          {
            loader: 'css-loader',
          },
          {
            loader: 'sass-loader',
            options: {
              sassOptions: {
                includePaths: ['node_modules/compass-mixins/lib', 'sass'],
              },
            },
          },
        ],
      },
    ],
  },
  externals: {
    react: 'React',
    redux: 'Redux',
    'react-redux': 'ReactRedux',
    'prop-types': 'PropTypes',
    'react-bootstrap': 'ReactBootstrap',
  },
  output: {
    devtoolNamespace: PLUGIN_ID,
    path: path.join(__dirname, '/dist'),
    publicPath: '/',
    filename: 'main.js',
  },
};