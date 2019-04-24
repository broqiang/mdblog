require('./bootstrap');

import hljs from './highlight';

    document.querySelectorAll('pre code').forEach((block) => {
        hljs.highlightBlock(block);
    });
