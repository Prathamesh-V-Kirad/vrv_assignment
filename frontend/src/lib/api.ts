import axios from 'axios';
  import { Task } from '@/types/task';
  
  const api = {
    getTasks: () => axios.get<Task[]>('http://localhost:3000/api/tasks').then(res => res.data),
    createTask: (task: Partial<Task>) => axios.post<Task>('http://localhost:3000/api/tasks', task).then(res => res.data),
    updateTask: (id: string, task: Partial<Task>) => axios.put<Task>(`http://localhost:3000/api/tasks/${id}`, task).then(res => res.data),
    deleteTask: (id: string) => axios.delete(`http://localhost:3000/api/tasks/${id}`).then(res => res.data)
  };
  
  export default api;