// 书籍表单页面脚本
document.addEventListener('DOMContentLoaded', function() {
    const bookForm = document.getElementById('book-form');
    const bookId = document.getElementById('book-id').value;
    
    // 如果有ID，说明是编辑模式，加载书籍数据
    if (bookId) {
        loadBookData(bookId);
    }
    
    // 表单提交事件
    bookForm.addEventListener('submit', function(e) {
        e.preventDefault();
        saveBook();
    });
});

// 加载书籍数据
async function loadBookData(bookId) {
    try {
        const book = await apiRequest(`/books/${bookId}`);
        fillFormWithBookData(book);
    } catch (error) {
        showError('加载书籍数据失败: ' + error.message);
    }
}

// 填充表单数据
function fillFormWithBookData(book) {
    document.getElementById('title').value = book.title;
    document.getElementById('author').value = book.author;
    document.getElementById('publisher').value = book.publisher || '';
    document.getElementById('publish_year').value = book.publish_year || '';
    document.getElementById('isbn').value = book.isbn;
    document.getElementById('description').value = book.description || '';
    document.getElementById('quantity').value = book.quantity;
}

// 从表单获取数据
function getBookDataFromForm() {
    return {
        title: document.getElementById('title').value,
        author: document.getElementById('author').value,
        publisher: document.getElementById('publisher').value,
        publish_year: parseInt(document.getElementById('publish_year').value) || null,
        isbn: document.getElementById('isbn').value,
        description: document.getElementById('description').value,
        quantity: parseInt(document.getElementById('quantity').value) || 0
    };
}

// 保存书籍
async function saveBook() {
    const bookId = document.getElementById('book-id').value;
    const bookData = getBookDataFromForm();
    
    try {
        let result;
        
        if (bookId) {
            // 更新现有书籍
            result = await apiRequest(`/books/${bookId}`, 'PUT', bookData);
            showSuccess('书籍已成功更新');
        } else {
            // 创建新书籍
            result = await apiRequest('/books', 'POST', bookData);
            showSuccess('书籍已成功添加');
        }
        
        // 跳转到书籍列表页面
        setTimeout(() => {
            window.location.href = '/books';
        }, 1000);
    } catch (error) {
        showError('保存书籍失败: ' + error.message);
    }
} 