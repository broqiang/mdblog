// 文档导航功能
class DocumentNavigation {
  constructor() {
    this.navContainer = document.getElementById("post-nav-list");
    this.contentContainer = document.getElementById("post-content");
    this.headings = [];
    this.currentActiveId = null;
    this.isScrolling = false;

    this.init();
  }

  init() {
    if (!this.navContainer || !this.contentContainer) {
      return;
    }

    this.extractHeadings();
    this.generateNavigation();
    this.setupScrollListener();
    this.updateActiveLink();
  }

  // 提取文档中的标题
  extractHeadings() {
    const headingElements = this.contentContainer.querySelectorAll(
      "h1, h2, h3, h4, h5, h6"
    );

    this.headings = Array.from(headingElements).map((heading) => {
      const level = parseInt(heading.tagName.charAt(1));
      let id = heading.id;

      // 如果没有ID，生成一个
      if (!id) {
        id = this.generateId(heading.textContent);
        heading.id = id;
      }

      return {
        id: id,
        text: heading.textContent.trim(),
        level: level,
        element: heading,
      };
    });
  }

  // 生成标题ID
  generateId(text) {
    return text
      .toLowerCase()
      .replace(/[^a-z0-9\u4e00-\u9fa5]+/g, "-")
      .replace(/^-+|-+$/g, "")
      .substring(0, 50);
  }

  // 生成导航HTML
  generateNavigation() {
    if (this.headings.length === 0) {
      // 如果没有标题，隐藏导航
      const navElement = document.getElementById("post-nav");
      if (navElement) {
        navElement.style.display = "none";
      }
      return;
    }

    const navHTML = this.headings
      .map((heading) => {
        return `
        <li class="post-nav-item">
          <a href="#${heading.id}" 
             class="post-nav-link level-${heading.level}" 
             data-target="${heading.id}">
            ${this.escapeHtml(heading.text)}
          </a>
        </li>
      `;
      })
      .join("");

    this.navContainer.innerHTML = navHTML;

    // 添加点击事件
    this.navContainer.addEventListener("click", this.handleNavClick.bind(this));
  }

  // HTML转义
  escapeHtml(text) {
    const div = document.createElement("div");
    div.textContent = text;
    return div.innerHTML;
  }

  // 处理导航链接点击
  handleNavClick(event) {
    event.preventDefault();

    const link = event.target.closest(".post-nav-link");
    if (!link) return;

    const targetId = link.getAttribute("data-target");
    const targetElement = document.getElementById(targetId);

    if (targetElement) {
      this.isScrolling = true;
      this.scrollToElement(targetElement);

      // 更新URL hash，但不触发滚动
      history.replaceState(null, null, `#${targetId}`);

      // 500ms后重新启用滚动监听
      setTimeout(() => {
        this.isScrolling = false;
        this.updateActiveLink();
      }, 500);
    }
  }

  // 平滑滚动到指定元素
  scrollToElement(element) {
    const headerHeight = 80; // header高度
    const extraOffset = 20; // 额外偏移
    const elementTop = element.offsetTop - headerHeight - extraOffset;

    window.scrollTo({
      top: elementTop,
      behavior: "smooth",
    });
  }

  // 设置滚动监听
  setupScrollListener() {
    let ticking = false;

    const scrollHandler = () => {
      if (!ticking && !this.isScrolling) {
        requestAnimationFrame(() => {
          this.updateActiveLink();
          ticking = false;
        });
        ticking = true;
      }
    };

    window.addEventListener("scroll", scrollHandler, { passive: true });

    // 页面加载时如果URL有hash，滚动到对应位置
    if (window.location.hash) {
      setTimeout(() => {
        const targetElement = document.querySelector(window.location.hash);
        if (targetElement) {
          this.scrollToElement(targetElement);
        }
      }, 100);
    }
  }

  // 更新当前激活的导航链接
  updateActiveLink() {
    if (this.headings.length === 0) return;

    const scrollTop = window.pageYOffset || document.documentElement.scrollTop;
    const headerHeight = 80;
    const offset = 100; // 提前激活的偏移量

    let activeId = null;

    // 从下往上查找第一个在视窗上方的标题
    for (let i = this.headings.length - 1; i >= 0; i--) {
      const heading = this.headings[i];
      const elementTop = heading.element.offsetTop - headerHeight - offset;

      if (scrollTop >= elementTop) {
        activeId = heading.id;
        break;
      }
    }

    // 如果没有找到，使用第一个标题
    if (!activeId && this.headings.length > 0) {
      activeId = this.headings[0].id;
    }

    // 更新激活状态
    if (activeId !== this.currentActiveId) {
      this.currentActiveId = activeId;

      // 移除所有激活状态
      this.navContainer.querySelectorAll(".post-nav-link").forEach((link) => {
        link.classList.remove("active");
      });

      // 添加新的激活状态
      if (activeId) {
        const activeLink = this.navContainer.querySelector(
          `[data-target="${activeId}"]`
        );
        if (activeLink) {
          activeLink.classList.add("active");
        }
      }
    }
  }
}

// 页面加载完成后初始化
document.addEventListener("DOMContentLoaded", () => {
  new DocumentNavigation();
});
