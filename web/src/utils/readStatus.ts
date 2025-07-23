interface ReadTopics {
  [topicId: string]: number; // topic ID -> highest read post ID
}

const STORAGE_KEY = 'minibb_read_topics';

export function getReadTopics(): ReadTopics {
  const stored = localStorage.getItem(STORAGE_KEY);
  if (!stored) return {};
  try {
    return JSON.parse(stored);
  } catch {
    return {};
  }
}

export function markTopicRead(topicId: number, highestPostId: number): void {
  const readTopics = getReadTopics();
  const topicIdStr = topicId.toString();
  
  // Only update if the new post ID is higher than the current one
  if (!readTopics[topicIdStr] || readTopics[topicIdStr] < highestPostId) {
    readTopics[topicIdStr] = highestPostId;
    localStorage.setItem(STORAGE_KEY, JSON.stringify(readTopics));
  }
}

export function hasUnreadPosts(topicId: number, lastPostId: number): boolean {
  const readTopics = getReadTopics();
  const topicIdStr = topicId.toString();
  const readUpTo = readTopics[topicIdStr] || 0;
  return lastPostId > readUpTo;
}

export function getUnreadCount(topicId: number, totalPosts: number): number {
  const readTopics = getReadTopics();
  const topicIdStr = topicId.toString();
  const readUpTo = readTopics[topicIdStr] || 0;
  
  // Find the number of posts that are read up to the stored post ID
  // This is an approximation since we don't have all post IDs
  return Math.max(0, totalPosts - Math.min(totalPosts, readUpTo));
}
