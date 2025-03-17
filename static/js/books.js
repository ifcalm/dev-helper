// 书籍列表页面脚本
document.addEventListener('DOMContentLoaded', function() {
    loadBooks();
    
    // 设置删除确认按钮事件
    const confirmDeleteBtn = document.getElementById('confirmDelete');
    confirmDeleteBtn.addEventListener('click', function() {
        const bookId = confirmDeleteBtn.dataset.bookId;
        deleteBook(bookId);
    });
});

// 加载所有书籍
async function loadBooks() {
    try {
        const books = await apiRequest('/books');
        displayBooks(books);
    } catch (error) {
        showError('加载书籍失败: ' + error.message);
    }
}

// 显示书籍列表
function displayBooks(books) {
    const tableBody = document.querySelector('#books-table tbody');
    tableBody.innerHTML = '';
    
    if (books.length === 0) {
        const row = document.createElement('tr');
        row.innerHTML = '<td colspan="8" class="text-center">没有找到书籍</td>';
        tableBody.appendChild(row);
        return;
    }
    
    books.forEach(book => {
        const row = document.createElement('tr');
        row.innerHTML = `
            <td>${book.ID}</td>
            <td>${book.title}</td>
            <td>${book.author}</td>
            <td>${book.publisher || '-'}</td>
            <td>${book.publish_year || '-'}</td>
            <td>${book.isbn}</td>
            <td>${book.quantity}</td>
            <td class="action-buttons">
                <a href="/books/edit/${book.ID}" class="btn btn-sm btn-primary">
                    <i class="bi bi-pencil"></i> 编辑
                </a>
                <button class="btn btn-sm btn-danger delete-book" data-book-id="${book.ID}">
                    <i class="bi bi-trash"></i> 删除
                </button>
            </td>
        `;
        tableBody.appendChild(row);
    });
    
    // 添加删除按钮事件
    document.querySelectorAll('.delete-book').forEach(button => {
        button.addEventListener('click', function() {
            const bookId = this.dataset.bookId;
            showDeleteConfirmation(bookId);
        });
    });
}

// 显示删除确认对话框
function showDeleteConfirmation(bookId) {
    const confirmDeleteBtn = document.getElementById('confirmDelete');
    confirmDeleteBtn.dataset.bookId = bookId;
    
    const deleteModal = new bootstrap.Modal(document.getElementById('deleteModal'));
    deleteModal.show();
}

// 删除书籍
async function deleteBook(bookId) {
    try {
        await apiRequest(`/books/${bookId}`, 'DELETE');
        
        // 关闭模态框
        const deleteModal = bootstrap.Modal.getInstance(document.getElementById('deleteModal'));
        deleteModal.hide();
        
        // 重新加载书籍列表
        loadBooks();
        showSuccess('书籍已成功删除');
    } catch (error) {
        showError('删除书籍失败: ' + error.message);
    }
} 