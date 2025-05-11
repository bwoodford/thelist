// API client for interacting with the backend
export interface Item {
  id: number;
  title: string;
  description: string;
  createdDate: string;
  modifiedDate?: string;
  completedDate?: string;
  isActive: boolean;
}

export interface NewItem {
  title: string;
  description: string;
  isActive: boolean;
}

const API_URL = 'http://localhost:8080'; // Adjust this to your API's URL

export async function getItems(): Promise<Item[]> {
  const response = await fetch(`${API_URL}/items`);
  if (!response.ok) {
    throw new Error('Failed to fetch items');
  }
  return response.json();
}

export async function createItem(item: NewItem): Promise<Item> {
  const response = await fetch(`${API_URL}/items`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(item),
  });
  
  if (!response.ok) {
    throw new Error('Failed to create item');
  }
  
  return response.json();
}

export async function deleteItem(id: number): Promise<void> {
  const response = await fetch(`${API_URL}/items/${id}`, {
    method: 'DELETE',
  });
  
  if (!response.ok) {
    throw new Error('Failed to delete item');
  }
}
