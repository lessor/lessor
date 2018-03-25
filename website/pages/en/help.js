/**
 * Copyright (c) 2017-present, Facebook, Inc.
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

const React = require('react');

const CompLibrary = require('../../core/CompLibrary.js');
const Container = CompLibrary.Container;
const GridBlock = CompLibrary.GridBlock;

const siteConfig = require(process.cwd() + '/siteConfig.js');

class Help extends React.Component {
  render() {
    const supportLinks = [
      {
        content:
          'Learn more using the [documentation on this site.](/docs/en/doc1.html)',
        title: 'Browse Docs',
      },
      {
        content: 'Ask questions about the documentation and project',
        title: 'Join the community',
      },
      {
        content: "Find out what's new with this project",
        title: 'Stay up to date',
      },
    ];

    return (
      <div className="docMainWrapper wrapper">
        <Container className="mainContainer documentContainer postContainer">
          <div className="post">
            <header className="postHeader">
              <h2>Where's the community?</h2>
            </header>
            <h2>Nulla</h2>

            <p>Nulla facilisi. Maecenas sodales nec purus eget posuere. Sed
      sapien quam, pretium a risus in, porttitor dapibus erat. Sed sit amet
      fringilla ipsum, eget iaculis augue. Integer sollicitudin tortor quis
      ultricies aliquam. Suspendisse fringilla nunc in tellus cursus, at
      placerat tellus scelerisque. Sed tempus elit a sollicitudin rhoncus.
      Nulla facilisi. Morbi nec dolor dolor. Orci varius natoque penatibus et
      magnis dis parturient montes, nascetur ridiculus mus. Cras et aliquet
      lectus. Pellentesque sit amet eros nisi. Quisque ac sapien in sapien
      congue accumsan. Nullam in posuere ante. Vestibulum ante ipsum primis in
      faucibus orci luctus et ultrices posuere cubilia Curae; Proin lacinia leo
      a nibh fringilla pharetra.</p>

            <h2>Orci</h2>

            <p>Orci varius natoque penatibus et magnis dis parturient montes,
      nascetur ridiculus mus. Proin venenatis lectus dui, vel ultrices ante
      bibendum hendrerit. Aenean egestas feugiat dui id hendrerit. Orci varius
      natoque penatibus et magnis dis parturient montes, nascetur ridiculus
      eget nisl sodales fermentum.</p>

            <GridBlock contents={supportLinks} layout="threeColumn" />
          </div>
        </Container>
      </div>
    );
  }
}

module.exports = Help;
