'use client';

import { useState } from 'react';
import ItemList from './components/ItemList';
import ItemForm from './components/ItemForm';

export default function Home() {
  const [showForm, setShowForm] = useState(false);
  const [refreshTrigger, setRefreshTrigger] = useState(0);

  const handleItemAdded = () => {
    // Trigger a refresh of the item list
    setRefreshTrigger(prev => prev + 1);
    // Optionally close the form after successful submission
    // setShowForm(false);
  };

  return (
    <div className="min-h-screen p-8 max-w-4xl mx-auto">
      <header className="mb-8 text-center">
        <h1 className="text-3xl font-bold mb-2">Item Manager</h1>
        <p className="text-gray-600 dark:text-gray-300">Manage your items with ease</p>
      </header>

      <div className="mb-6 flex justify-end">
        <button
          onClick={() => setShowForm(!showForm)}
          className="bg-blue-600 hover:bg-blue-700 text-white font-medium py-2 px-4 rounded-md transition-colors"
        >
          {showForm ? 'Hide Form' : 'Add New Item'}
        </button>
      </div>

      {showForm && (
        <div className="mb-8">
          <ItemForm onItemAdded={handleItemAdded} />
        </div>
      )}

      <div key={refreshTrigger}>
        <ItemList />
      </div>
    </div>
  );
}
