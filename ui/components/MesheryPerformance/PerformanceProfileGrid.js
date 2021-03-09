//@ts-check
import React from 'react';
import PerformanceCard from './PerformanceCard';

/**
 * PerformanceProfileGrid is the react component for rendering
 * grid
 * @param {{
 *  profiles: Array<{
 *    id: string,
 *    created_at: string,
 *    updated_at: string,
 *    endpoints: Array<string>,
 *    load_generators: Array<string>,
 *    name: string,
 *    user_id: string,
 *    duration: string,
 *  }>
 * }} props props
 */
function PerformanceProfileGrid({ profiles }) {
  return (
    <div>
      {profiles.map(() => <PerformanceCard />)}
    </div>
  )
}

export default PerformanceProfileGrid
