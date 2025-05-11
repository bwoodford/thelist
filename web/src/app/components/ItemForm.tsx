'use client';

import { useState } from 'react';
import { NewItem, createItem } from '../api/itemsApi';

interface ItemFormProps {
  onItemAdded: () => void;
}

export default function ItemForm({ onItemAdded }: ItemFormProps) {
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!title.trim()) {
      setError('Title is required');
      return;
    }
    
    const newItem: NewItem = {
      title,
      description,
      isActive: true
    };
    
    try {
      setSubmitting(true);
      await createItem(newItem);
      setTitle('');
      setDescription('');
      setError(null);
      onItemAdded();
    } catch (err) {
      setError('Failed to create item. Please try again.');
      console.error(err);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="w-full bg-white dark:bg-gray-800 p-6 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700">
      <h2 className="text-xl font-bold mb-4">Add New Item</h2>
      
      {error && (
        <div className="bg-red-100 dark:bg-red-900/30 text-red-700 dark:text-red-400 p-3 rounded mb-4">
          {error}
        </div>
      )}
      
      <form onSubmit={handleSubmit}>
        <div className="mb-4">
          <label htmlFor="title" className="block text-sm font-medium mb-1">
            Title *
          </label>
          <input
            type="text"
            id="title"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md 
                      dark:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="Enter title"
            disabled={submitting}
          />
        </div>
        
        <div className="mb-4">
          <label htmlFor="description" className="block text-sm font-medium mb-1">
            Description
          </label>
          <textarea
            id="description"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md 
                      dark:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="Enter description"
            rows={3}
            disabled={submitting}
          />
        </div>
        
        <button
          type="submit"
          disabled={submitting}
          className="w-full bg-blue-600 hover:bg-blue-700 text-white font-medium py-2 px-4 
                    rounded-md transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {submitting ? 'Adding...' : 'Add Item'}
        </button>
      </form>
    </div>
  );
}
