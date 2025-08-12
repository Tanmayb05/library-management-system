import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { Plus, Edit, Trash2, Book, Search, X } from 'lucide-react';
import './App.css';

const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080/api/v1';

function App() {
  const [books, setBooks] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [showModal, setShowModal] = useState(false);
  const [editingBook, setEditingBook] = useState(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [formData, setFormData] = useState({
    title: '',
    author: '',
    isbn: '',
    publication_year: '',
    available: true
  });

  useEffect(() => {
    fetchBooks();
  }, []);

  const fetchBooks = async () => {
    try {
      setLoading(true);
      const response = await axios.get(`${API_URL}/books`);
      setBooks(response.data.data || []);
      setError('');
    } catch (err) {
      setError('Failed to fetch books');
      console.error('Error fetching books:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      if (editingBook) {
        // Update existing book
        const updateData = {};
        if (formData.title !== editingBook.title) updateData.title = formData.title;
        if (formData.author !== editingBook.author) updateData.author = formData.author;
        if (formData.isbn !== editingBook.isbn) updateData.isbn = formData.isbn;
        if (formData.publication_year !== editingBook.publication_year) {
          updateData.publication_year = parseInt(formData.publication_year);
        }
        if (formData.available !== editingBook.available) updateData.available = formData.available;

        await axios.put(`${API_URL}/books/${editingBook.id}`, updateData);
      } else {
        // Create new book
        const bookData = {
          ...formData,
          publication_year: parseInt(formData.publication_year) || 0
        };
        await axios.post(`${API_URL}/books`, bookData);
      }
      
      fetchBooks();
      handleCloseModal();
      setError('');
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to save book');
    }
  };

  const handleDelete = async (id) => {
    if (window.confirm('Are you sure you want to delete this book?')) {
      try {
        await axios.delete(`${API_URL}/books/${id}`);
        fetchBooks();
        setError('');
      } catch (err) {
        setError('Failed to delete book');
      }
    }
  };

  const handleEdit = (book) => {
    setEditingBook(book);
    setFormData({
      title: book.title,
      author: book.author,
      isbn: book.isbn,
      publication_year: book.publication_year.toString(),
      available: book.available
    });
    setShowModal(true);
  };

  const handleCloseModal = () => {
    setShowModal(false);
    setEditingBook(null);
    setFormData({
      title: '',
      author: '',
      isbn: '',
      publication_year: '',
      available: true
    });
  };

  const filteredBooks = books.filter(book =>
    book.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
    book.author.toLowerCase().includes(searchTerm.toLowerCase()) ||
    book.isbn.includes(searchTerm)
  );

  if (loading) {
    return (
      <div className="app">
        <div className="loading">
          <div className="spinner"></div>
          <p>Loading books...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="app">
      <header className="header">
        <div className="header-content">
          <div className="header-title">
            <Book className="header-icon" />
            <h1>Library Management System</h1>
          </div>
          <button 
            className="btn btn-primary"
            onClick={() => setShowModal(true)}
          >
            <Plus size={20} />
            Add Book
          </button>
        </div>
      </header>

      <main className="main">
        <div className="search-section">
          <div className="search-box">
            <Search className="search-icon" />
            <input
              type="text"
              placeholder="Search books by title, author, or ISBN..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="search-input"
            />
          </div>
        </div>

        {error && (
          <div className="error-message">
            {error}
          </div>
        )}

        <div className="books-grid">
          {filteredBooks.length === 0 ? (
            <div className="empty-state">
              <Book size={64} />
              <h3>No books found</h3>
              <p>
                {searchTerm ? 'Try adjusting your search terms' : 'Add your first book to get started'}
              </p>
            </div>
          ) : (
            filteredBooks.map((book) => (
              <div key={book.id} className="book-card">
                <div className="book-header">
                  <h3 className="book-title">{book.title}</h3>
                  <div className="book-actions">
                    <button
                      className="btn-icon"
                      onClick={() => handleEdit(book)}
                      title="Edit book"
                    >
                      <Edit size={16} />
                    </button>
                    <button
                      className="btn-icon btn-danger"
                      onClick={() => handleDelete(book.id)}
                      title="Delete book"
                    >
                      <Trash2 size={16} />
                    </button>
                  </div>
                </div>
                
                <div className="book-details">
                  <p><strong>Author:</strong> {book.author}</p>
                  <p><strong>ISBN:</strong> {book.isbn}</p>
                  <p><strong>Year:</strong> {book.publication_year}</p>
                  <div className="book-status">
                    <span className={`status ${book.available ? 'available' : 'unavailable'}`}>
                      {book.available ? 'Available' : 'Unavailable'}
                    </span>
                  </div>
                </div>
              </div>
            ))
          )}
        </div>
      </main>

      {showModal && (
        <div className="modal-overlay">
          <div className="modal">
            <div className="modal-header">
              <h2>{editingBook ? 'Edit Book' : 'Add New Book'}</h2>
              <button 
                className="btn-icon"
                onClick={handleCloseModal}
              >
                <X size={20} />
              </button>
            </div>

            <form onSubmit={handleSubmit} className="modal-form">
              <div className="form-group">
                <label htmlFor="title">Title *</label>
                <input
                  type="text"
                  id="title"
                  value={formData.title}
                  onChange={(e) => setFormData({...formData, title: e.target.value})}
                  required
                />
              </div>

              <div className="form-group">
                <label htmlFor="author">Author *</label>
                <input
                  type="text"
                  id="author"
                  value={formData.author}
                  onChange={(e) => setFormData({...formData, author: e.target.value})}
                  required
                />
              </div>

              <div className="form-group">
                <label htmlFor="isbn">ISBN *</label>
                <input
                  type="text"
                  id="isbn"
                  value={formData.isbn}
                  onChange={(e) => setFormData({...formData, isbn: e.target.value})}
                  required
                />
              </div>

              <div className="form-group">
                <label htmlFor="publication_year">Publication Year</label>
                <input
                  type="number"
                  id="publication_year"
                  value={formData.publication_year}
                  onChange={(e) => setFormData({...formData, publication_year: e.target.value})}
                  min="1000"
                  max={new Date().getFullYear() + 1}
                />
              </div>

              <div className="form-group">
                <label className="checkbox-label">
                  <input
                    type="checkbox"
                    checked={formData.available}
                    onChange={(e) => setFormData({...formData, available: e.target.checked})}
                  />
                  Available
                </label>
              </div>

              <div className="modal-actions">
                <button 
                  type="button" 
                  className="btn btn-secondary"
                  onClick={handleCloseModal}
                >
                  Cancel
                </button>
                <button 
                  type="submit" 
                  className="btn btn-primary"
                >
                  {editingBook ? 'Update' : 'Create'} Book
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}

export default App;