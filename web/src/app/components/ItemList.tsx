'use client';

import { useState, useEffect } from 'react';
import { Item, getItems, deleteItem } from '../api/itemsApi';

export default function ItemList() {
  const [items, setItems] = useState<Item[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchItems();
  }, []);

  const fetchItems = async () => {
    try {
      setLoading(true);
      const data = await getItems();
      // Sort items by created date (newest first)
      const sortedItems = data.sort((a, b) => 
        new Date(b.createdDate).getTime() - new Date(a.createdDate).getTime()
      );
      setItems(sortedItems);
      setError(null);
    } catch (err) {
      setError('Failed to load items. Please try again later.');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (id: number) => {
    try {
      await deleteItem(id);
      // Remove the item from the local state
      setItems(items.filter(item => item.id !== id));
    } catch (err) {
      setError('Failed to delete item');
      console.error(err);
    }
  };

  if (loading) return <div className="text-center py-4">Loading items...</div>;
  if (error) return <div className="text-center py-4 text-red-500">{error}</div>;
  if (items.length === 0) return <div className="text-center py-4">No items found. Add some!</div>;

  return (
    <div className="w-full">
      <h2 className="text-xl font-bold mb-4">Your Items</h2>
      <ul className="space-y-4">
        {items.map((item) => (
          <li key={item.id} className="border border-gray-200 dark:border-gray-700 rounded-lg p-4 bg-white dark:bg-gray-800">
            <div className="flex justify-between items-start">
              <div>
                <h3 className="font-semibold text-lg">{item.title}</h3>
                <p className="text-gray-600 dark:text-gray-300 mt-1">{item.description}</p>
                <div className="text-xs text-gray-500 dark:text-gray-400 mt-2">
                  Created: {new Date(item.createdDate).toLocaleString()}
                </div>
              </div>
              <button 
                onClick={() => handleDelete(item.id)}
                className="text-red-500 hover:text-red-700 text-sm"
              >
                Delete
              </button>
            </div>
          </li>
        ))}
      </ul>
    </div>
  );
}
