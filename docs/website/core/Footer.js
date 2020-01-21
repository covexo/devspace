/**
 * Copyright (c) 2017-present, Facebook, Inc.
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

const React = require("react");

class Footer extends React.Component {
  docUrl(doc, language) {
    const baseUrl = this.props.config.baseUrl;
    const docsUrl = this.props.config.docsUrl;
    const docsPart = `${docsUrl ? `${docsUrl}/` : ""}`;
    const langPart = `${language ? `${language}/` : ""}`;
    return `${baseUrl}${docsPart}${langPart}${doc}`;
  }

  pageUrl(doc, language) {
    const baseUrl = this.props.config.baseUrl;
    return baseUrl + (language ? `${language}/` : "") + doc;
  }

  render() {
    let optionalComponent;

    try {
      const Chat = require("@devspace/react-components").Chat;
      const Analytics = require("@devspace/react-components").Analytics;

      optionalComponent = (
        <div>
          <Chat />
          <Analytics />
        </div>
      );
    } catch (e) {}

    return (
      <footer className="nav-footer" id="footer">

        <div className="star-button">
          <script async defer src="https://buttons.github.io/buttons.js"></script>
          <a className="github-button" href="https://github.com/devspace-cloud/devspace" data-size="large" data-show-count="true" aria-label="Star devspace-cloud/devspace on GitHub">Star</a>
        </div>
        
        <script type="text/javascript" dangerouslySetInnerHTML={{__html: `
        var versionMeta = document.querySelector("head > meta[name='docsearch:version']");
        var sidebarVersions = {
          "v3.5.18": "v3", 
          "v4.0.0": "v4", 
          "v4.0.3": "v4",
          "v4.1.0": "v4",
          "v4.2.0": "v4.2",
          "v4.3.5": "v4.2",
        };

        if (versionMeta) {
          let version = versionMeta.getAttribute("content");
          let sidebarVersionsArray = Object.keys(sidebarVersions);
          let sidebarVersion = sidebarVersions[sidebarVersionsArray[sidebarVersionsArray.length - 1]];
          
          if (version != "next") {
            let versionSplit = version.split(".");
            let major = versionSplit[0].substr(1);
            let minor = versionSplit[1];
            let revision = versionSplit[2];

            for (let versionKey in sidebarVersions) {
              let sidebarVersionSplit = versionKey.split(".");
              let sidebarMajor = sidebarVersionSplit[0].substr(1);
              let sidebarMinor = sidebarVersionSplit[1];
              let sidebarRevision = sidebarVersionSplit[2];

              if (major > sidebarMajor || (major == sidebarMajor && minor > sidebarMinor) || (major == sidebarMajor && minor == sidebarMinor && revision >= sidebarRevision)) {
                sidebarVersion = sidebarVersions[versionKey];
              } else {
                break;
              }
            }
          }

          document.querySelector("body").setAttribute("data-version", version);
          document.querySelector("body").setAttribute("data-sidebar-version", sidebarVersion);
        }

        if (location.hostname == "devspace.cloud") {
          document.querySelector(".headerWrapper > header > a:nth-child(2)").setAttribute("href", "/docs/versions");
        }
        `}}>
        </script>

        {optionalComponent}

        <noscript>
          <iframe
            src="https://www.googletagmanager.com/ns.html?id=GTM-KM6KSWG"
            height="0"
            width="0"
            style={{ display: "none", visibility: "hidden" }}
          />
        </noscript>
      </footer>
    );
  }
}

module.exports = Footer;
