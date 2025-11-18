function createVarReader(root = document.documentElement) {
  const cs = getComputedStyle(root);
  const cache = new Map();
  return function v(name, fallback = '') {
    if (cache.has(name)) return cache.get(name);
    const val = cs.getPropertyValue(name).trim() || fallback;
    cache.set(name, val);
    return val;
  };
}

export function makeTheme(root) {
  const v = createVarReader(root);

  return {
    colors: {
      primary: {
        main: v('--colors-primary-main', '#32329f'),
        light: v('--colors-primary-light', '#8e8edc'),
        dark: v('--colors-primary-dark', '#0d0d2b'),
        contrastText: v('--colors-primary-contrastText', '#ffffff'),
      },
      success: {
        main: v('--colors-success-main', '#00aa13'),
        light: v('--colors-success-light', '#44ff59'),
        dark: v('--colors-success-dark', '#001102'),
        contrastText: v('--colors-success-contrastText', '#000000'),
      },
      warning: {
        main: v('--colors-warning-main', '#d4ad03'),
        light: v('--colors-warning-light', '#fde373'),
        dark: v('--colors-warning-dark', '#3d3201'),
        contrastText: v('--colors-warning-contrastText', '#ffffff'),
      },
      error: {
        main: v('--colors-error-main', '#e53935'),
        light: v('--colors-error-light', '#f6bebd'),
        dark: v('--colors-error-dark', '#72110f'),
        contrastText: v('--colors-error-contrastText', '#000000'),
      },
      text: {
        primary: v('--colors-text-primary', '#333333'),
        secondary: v('--colors-text-secondary', '#808080'),
      },
      border: {
        dark: v('--colors-border-dark', 'rgba(0,0,0, 0.1)'),
        light: v('--colors-border-light', '#ffffff'),
      },
      responses: {
        success: {
          color: v('--colors-responses-success-color', '#00aa13'),
          backgroundColor: v('--colors-responses-success-background-color', 'rgba(0,170,19,0.1)'),
        },
        error: {
          color: v('--colors-responses-error-color', '#e53935'),
          backgroundColor: v('--colors-responses-error-background-color', 'rgba(229,57,53,0.1)'),
        },
        redirect: {
          color: v('--colors-responses-redirect-color', '#ffa500'),
          backgroundColor: v('--colors-responses-redirect-background-color', 'rgba(255,165,0,0.1)'),
        },
        info: {
          color: v('--colors-responses-info-color', '#87ceeb'),
          backgroundColor: v('--colors-responses-info-background-color', 'rgba(135,206,235,0.1)'),
        },
      },
      http: {
        get: v('--colors-http-get', '#6bbd5b'),
        post: v('--colors-http-post', '#248fb2'),
        put: v('--colors-http-put', '#9b708b'),
        options: v('--colors-http-options', '#d3ca12'),
        patch: v('--colors-http-patch', '#e09d43'),
        delete: v('--colors-http-delete', '#e27a7a'),
        basic: v('--colors-http-basic', '#999999'),
        link: v('--colors-http-link', '#31bbb6'),
        head: v('--colors-http-head', '#c167e4'),
      },
      gray: {
        50: v('--colors-gray-50', '#fafafa'),
        100: v('--colors-gray-100', '#f5f5f5'),
      }
    },
    schema: {
      linesColor: v('--schema-lines-color', '#a4a4c6'),
      typeNameColor: v('--schema-type-name-color', '#808080'),
      typeTitleColor: v('--schema-type-title-color', '#808080'),
      requireLabelColor: v('--schema-require-label-color', '#e53935'),
      nestedBackground: v('--schema-nested-background', '#fafafa'),
      arrow: {
        color: v('--schema-arrow-color', '#808080'),
      },
    },
    typography: {
      code: {
        color: v('--typography-code-color', '#e53935'),
        backgroundColor: v('--typography-code-background-color', 'rgba(38, 50, 56, 0.05)'),
      },
      links: {
        color: v('--typography-links-color', '#32329f'),
        visited: v('--typography-links-visited', '#32329f'),
        hover: v('--typography-links-hover', '#6868cf'),
      },
    },
    sidebar: {
      backgroundColor: v('--sidebar-background-color', '#fafafa'),
      textColor: v('--sidebar-text-color', '#333333'),
      arrow: {
        color: v('--sidebar-arrow-color', '#333333'),
      },
    },
    rightPanel: {
      backgroundColor: v('--right-panel-background-color', '#263238'),
      textColor: v('--right-panel-text-color', '#ffffff'),
      servers: {
        overlay: {
          backgroundColor: v('--right-panel-servers-overlay-background-color', '#fafafa'),
          textColor: '#263238',
        },
        url: {
          backgroundColor: v('--right-panel-servers-url-background-color', '#ffffff'),
        }
      },
    },
    codeBlock: {
      backgroundColor: v('--code-block-background-color', '#11171a'),
    },
    fab: {
      backgroundColor: v('--fab-background-color', '#f2f2f2'),
      color: v('--fab-color', '#0065FB'),
    },
    extensionsHook: (c) => {
      if (c === 'UnderlinedHeader') {
        return {
          color:       v('--underlined-header-color',         '#a1a1aa'),
          fontWeight:  v('--underlined-header-font-weight',   'bold'),
          borderBottom:v('--underlined-header-border-bottom', '1px solid #3f3f46'),
        };
      }
    },
  };
}

const theme = makeTheme();
export default theme;
