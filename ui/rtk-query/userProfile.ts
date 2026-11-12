type UserProfileSummaryResponse = {
  id?: string;
  email?: string;
  userId?: string;
  user_id?: string;
  avatarUrl?: string;
  avatar_url?: string;
  firstName?: string;
  first_name?: string;
  lastName?: string;
  last_name?: string;
};

export const normalizeUserProfileSummary = (response?: UserProfileSummaryResponse) => {
  if (!response || typeof response !== 'object') {
    return undefined;
  }

  return {
    id: response.id ?? response.userId ?? response.user_id,
    email: response.email,
    userId: response.userId ?? response.user_id ?? response.id,
    avatarUrl: response.avatarUrl ?? response.avatar_url,
    firstName: response.firstName ?? response.first_name,
    lastName: response.lastName ?? response.last_name,
  };
};
