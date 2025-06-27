// 搜索功能
class SearchManager {
  constructor() {
    this.searchToggle = document.getElementById("search-toggle");
    this.searchModal = document.getElementById("search-modal");
    this.searchInput = document.getElementById("search-input");
    this.searchClose = document.getElementById("search-close");
    this.searchResults = document.getElementById("search-results");

    this.searchTimeout = null;
    this.isSearching = false;

    // 快捷键相关
    this.isMac = /Mac|iPod|iPhone|iPad/.test(navigator.platform);
    this.keyPressCount = 0;
    this.keyPressTimer = null;
    this.lastKeyTime = 0;
    this.doubleClickInterval = 300; // 300ms内双击有效

    this.init();
  }

  init() {
    this.setupEventListeners();
    this.setupKeyboardShortcuts();
  }

  setupEventListeners() {
    // 打开搜索弹窗
    if (this.searchToggle) {
      this.searchToggle.addEventListener("click", () => {
        this.openModal();
      });
    }

    // 关闭搜索弹窗
    if (this.searchClose) {
      this.searchClose.addEventListener("click", () => {
        this.closeModal();
      });
    }

    // 点击弹窗背景关闭
    if (this.searchModal) {
      this.searchModal.addEventListener("click", (e) => {
        if (e.target === this.searchModal) {
          this.closeModal();
        }
      });
    }

    // ESC键关闭弹窗
    document.addEventListener("keydown", (e) => {
      if (e.key === "Escape" && this.searchModal.classList.contains("active")) {
        this.closeModal();
      }
    });

    // 搜索输入事件
    if (this.searchInput) {
      this.searchInput.addEventListener("input", (e) => {
        this.handleSearchInput(e.target.value);
      });

      // 防止回车键提交表单
      this.searchInput.addEventListener("keydown", (e) => {
        if (e.key === "Enter") {
          e.preventDefault();
        }
      });
    }
  }

  setupKeyboardShortcuts() {
    let isTargetKeyPressed = false;

    document.addEventListener("keydown", (e) => {
      // 检查是否是目标键（Mac的Command或Windows的Ctrl）
      const isTargetKey = this.isMac
        ? e.metaKey && e.key === "Meta"
        : e.ctrlKey && e.key === "Control";

      if (isTargetKey && !isTargetKeyPressed) {
        isTargetKeyPressed = true;
        this.handleKeyPress();
      }
    });

    document.addEventListener("keyup", (e) => {
      // 重置按键状态
      if (
        (this.isMac && e.key === "Meta") ||
        (!this.isMac && e.key === "Control")
      ) {
        isTargetKeyPressed = false;
      }
    });

    // 监听窗口失焦，重置状态
    window.addEventListener("blur", () => {
      isTargetKeyPressed = false;
      this.resetKeyPress();
    });
  }

  handleKeyPress() {
    const currentTime = Date.now();

    // 如果距离上次按键时间超过间隔，重置计数
    if (currentTime - this.lastKeyTime > this.doubleClickInterval) {
      this.keyPressCount = 0;
    }

    this.keyPressCount++;
    this.lastKeyTime = currentTime;

    // 清除之前的定时器
    if (this.keyPressTimer) {
      clearTimeout(this.keyPressTimer);
    }

    // 如果是双击
    if (this.keyPressCount === 2) {
      this.openModal();
      this.resetKeyPress();
      return;
    }

    // 设置定时器，超时后重置计数
    this.keyPressTimer = setTimeout(() => {
      this.resetKeyPress();
    }, this.doubleClickInterval);
  }

  resetKeyPress() {
    this.keyPressCount = 0;
    this.lastKeyTime = 0;
    if (this.keyPressTimer) {
      clearTimeout(this.keyPressTimer);
      this.keyPressTimer = null;
    }
  }

  openModal() {
    if (this.searchModal) {
      this.searchModal.classList.add("active");
      // 延迟聚焦，确保弹窗动画完成
      setTimeout(() => {
        if (this.searchInput) {
          this.searchInput.focus();
        }
      }, 100);

      // 显示快捷键提示（仅首次使用时）
      this.showShortcutHint();
    }
  }

  showShortcutHint() {
    // 检查是否已经显示过提示
    if (localStorage.getItem("searchShortcutHintShown")) {
      return;
    }

    const keyName = this.isMac ? "Command" : "Ctrl";
    const hint = document.createElement("div");
    hint.className = "search-shortcut-hint";
    hint.innerHTML = `
      <div class="search-hint-content">
        <span>💡 小贴士：双击 ${keyName} 键可快速打开搜索</span>
        <button class="search-hint-close">&times;</button>
      </div>
    `;

    // 添加样式
    hint.style.cssText = `
      position: fixed;
      top: 20px;
      right: 20px;
      background: #333;
      color: white;
      padding: 12px 16px;
      border-radius: 8px;
      font-size: 14px;
      z-index: 1001;
      box-shadow: 0 4px 12px rgba(0,0,0,0.2);
      animation: slideInRight 0.3s ease;
    `;

    document.body.appendChild(hint);

    // 关闭按钮事件
    const closeBtn = hint.querySelector(".search-hint-close");
    closeBtn.addEventListener("click", () => {
      hint.remove();
    });

    // 3秒后自动消失
    setTimeout(() => {
      if (hint.parentNode) {
        hint.remove();
      }
    }, 3000);

    // 标记已显示过提示
    localStorage.setItem("searchShortcutHintShown", "true");
  }

  closeModal() {
    if (this.searchModal) {
      this.searchModal.classList.remove("active");
      this.clearResults();
      if (this.searchInput) {
        this.searchInput.value = "";
      }
    }
  }

  handleSearchInput(query) {
    // 清除之前的定时器
    if (this.searchTimeout) {
      clearTimeout(this.searchTimeout);
    }

    const trimmedQuery = query.trim();

    // 如果查询长度小于2，清空结果
    if (trimmedQuery.length < 2) {
      this.clearResults();
      return;
    }

    // 防抖：300ms后执行搜索
    this.searchTimeout = setTimeout(() => {
      this.performSearch(trimmedQuery);
    }, 300);
  }

  async performSearch(query) {
    if (this.isSearching) return;

    this.isSearching = true;
    this.showLoading();

    try {
      const response = await fetch(
        `/api/search?q=${encodeURIComponent(query)}&page=1&size=10`
      );

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      this.displayResults(data.posts, query);
    } catch (error) {
      console.error("搜索失败:", error);
      this.showError("搜索失败，请稍后重试");
    } finally {
      this.isSearching = false;
    }
  }

  showLoading() {
    if (this.searchResults) {
      this.searchResults.innerHTML =
        '<div class="search-loading">搜索中...</div>';
    }
  }

  showError(message) {
    if (this.searchResults) {
      this.searchResults.innerHTML = `<div class="search-no-results">${message}</div>`;
    }
  }

  clearResults() {
    if (this.searchResults) {
      this.searchResults.innerHTML = "";
    }
  }

  displayResults(posts, query) {
    if (!this.searchResults) return;

    if (!posts || posts.length === 0) {
      this.searchResults.innerHTML =
        '<div class="search-no-results">未找到相关文章</div>';
      return;
    }

    const resultsHTML = posts
      .map((post) => {
        const title = this.highlightKeyword(post.title, query);
        // 修复时间格式问题，使用updated_at字段
        const updateTime = post.updated_at
          ? new Date(post.updated_at)
          : new Date(post.created_at);
        const date = updateTime.toLocaleDateString("zh-CN", {
          year: "numeric",
          month: "2-digit",
          day: "2-digit",
        });

        return `
        <div class="search-result-item" data-id="${post.id}">
          <div class="search-result-title">${title}</div>
          <div class="search-result-meta">
            <span class="search-result-category">${
              post.category || "未分类"
            }</span>
            <span class="search-result-date">${date}</span>
          </div>
        </div>
      `;
      })
      .join("");

    this.searchResults.innerHTML = resultsHTML;

    // 添加点击事件
    this.searchResults
      .querySelectorAll(".search-result-item")
      .forEach((item) => {
        item.addEventListener("click", () => {
          const postId = item.dataset.id;
          if (postId) {
            window.location.href = `/post/${postId}`;
          }
        });
      });
  }

  highlightKeyword(text, keyword) {
    if (!keyword || !text) return text;

    const regex = new RegExp(`(${keyword})`, "gi");
    return text.replace(regex, '<strong style="color: #0056b3;">$1</strong>');
  }
}

// 分页功能
class PaginationManager {
  constructor() {
    this.currentPage = 1;
    this.pageSize = 10;
    this.total = 0;
  }

  setPagination(page, pageSize, total) {
    this.currentPage = page;
    this.pageSize = pageSize;
    this.total = total;
  }

  getTotalPages() {
    return Math.ceil(this.total / this.pageSize);
  }

  hasNextPage() {
    return this.currentPage < this.getTotalPages();
  }

  hasPrevPage() {
    return this.currentPage > 1;
  }

  getNextPage() {
    return this.hasNextPage() ? this.currentPage + 1 : this.currentPage;
  }

  getPrevPage() {
    return this.hasPrevPage() ? this.currentPage - 1 : this.currentPage;
  }
}

// 移动端菜单管理器
class MobileMenuManager {
  constructor() {
    this.menuToggle = document.getElementById("mobile-menu-toggle");
    this.mobileNav = document.getElementById("mobile-nav");
    this.isOpen = false;

    this.init();
  }

  init() {
    if (this.menuToggle && this.mobileNav) {
      this.setupEventListeners();
    }
  }

  setupEventListeners() {
    // 菜单按钮点击事件
    this.menuToggle.addEventListener("click", () => {
      this.toggleMenu();
    });

    // 点击菜单项后关闭菜单
    this.mobileNav.addEventListener("click", (e) => {
      if (e.target.tagName === "A") {
        this.closeMenu();
      }
    });

    // 点击页面其他区域关闭菜单
    document.addEventListener("click", (e) => {
      if (
        this.isOpen &&
        !this.menuToggle.contains(e.target) &&
        !this.mobileNav.contains(e.target)
      ) {
        this.closeMenu();
      }
    });

    // ESC键关闭菜单
    document.addEventListener("keydown", (e) => {
      if (e.key === "Escape" && this.isOpen) {
        this.closeMenu();
      }
    });

    // 窗口大小改变时关闭菜单
    window.addEventListener("resize", () => {
      if (window.innerWidth > 768 && this.isOpen) {
        this.closeMenu();
      }
    });
  }

  toggleMenu() {
    if (this.isOpen) {
      this.closeMenu();
    } else {
      this.openMenu();
    }
  }

  openMenu() {
    this.isOpen = true;
    this.menuToggle.classList.add("active");
    this.mobileNav.classList.add("active");
  }

  closeMenu() {
    this.isOpen = false;
    this.menuToggle.classList.remove("active");
    this.mobileNav.classList.remove("active");
  }
}

// 初始化管理器
document.addEventListener("DOMContentLoaded", () => {
  new SearchManager();
  new MobileMenuManager();
});
