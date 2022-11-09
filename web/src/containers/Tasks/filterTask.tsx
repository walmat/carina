import { Task } from '../../stores/Main/reducers/tasks';

export const filterTask = (task: Task, search: string) => {
  const safeSearch = search.trim().toLowerCase();

  const matches = (value: string) => value.toLowerCase().includes(safeSearch);

  if (matches(task.id)) {
    return true;
  }

  if (matches(task.store.name)) {
    return true;
  }

  if (matches(task.profile.name)) {
    return true;
  }

  if (matches(task.product.name)) {
    return true;
  }

  return matches(task.status);
}
