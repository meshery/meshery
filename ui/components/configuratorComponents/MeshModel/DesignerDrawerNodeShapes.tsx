import React from 'react';

export function TriangleShape({ color, styles }) {
  return (
    <svg width="40" height="40" viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
      <g clipPath="url(#clip0_18568_8277)">
        <g clipPath="url(#clip1_18568_8277)">
          <path
            d="M20.0137 0L30.6529 20.0141L40 40H20.0137H0L9.37457 20.0141L20.0137 0Z"
            fill={color}
            {...styles}
          />
        </g>
      </g>
      <defs>
        <clipPath id="clip0_18568_8277">
          <rect width="40" height="40" fill="none" />
        </clipPath>
        <clipPath id="clip1_18568_8277">
          <rect width="40" height="40" fill="none" />
        </clipPath>
      </defs>
    </svg>
  );
}

export function RoundTriangleShape({ color, styles }) {
  return (
    <svg width="40" height="40" viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
      <path
        d="M26.0622 3.5C23.3679 -1.16667 16.6321 -1.16667 13.9378 3.5L1.06218 25.8013C-1.63212 30.4679 1.73576 36.3013 7.12436 36.3013H32.8756C38.2642 36.3013 41.6321 30.4679 38.9378 25.8013L26.0622 3.5Z"
        fill={color}
        {...styles}
      />
    </svg>
  );
}

export function PentagonShape({ color, styles }) {
  return (
    <svg width="40" height="40" viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
      <g clipPath="url(#clip0_18568_8232)">
        <g clipPath="url(#clip1_18568_8232)">
          <path
            d="M19.9862 0L0 15.2712L7.62958 40H32.3428L40 15.2712L19.9862 0Z"
            fill={color}
            {...styles}
          />
        </g>
      </g>
      <defs>
        <clipPath id="clip0_18568_8232">
          <rect width="40" height="40" fill="none" />
        </clipPath>
        <clipPath id="clip1_18568_8232">
          <rect width="40" height="40" fill="none" />
        </clipPath>
      </defs>
    </svg>
  );
}

export function RoundPentagonShape({ color, styles }) {
  return (
    <svg width="40" height="40" viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
      <path
        d="M17.6249 0.76393C19.0411 -0.254644 20.9589 -0.254642 22.3751 0.76393L38.3341 12.2416C39.7503 13.2602 40.3429 15.0657 39.8019 16.7138L33.7063 35.2851C33.1653 36.9332 31.6137 38.049 29.8631 38.049H10.1369C8.38627 38.049 6.83477 36.9332 6.29382 35.2851L0.198047 16.7138C-0.342914 15.0657 0.249705 13.2602 1.66596 12.2416L17.6249 0.76393Z"
        fill={color}
        {...styles}
      />
    </svg>
  );
}
export function RoundRectangleShape({ color, styles }) {
  return (
    <svg width="40" height="40" viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
      <rect width="40" height="40" rx="9" fill={color} {...styles} />
    </svg>
  );
}
export function RectangleShape({ color, styles }) {
  return (
    <svg width="40" height="40" viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
      <rect width="40" height="40" fill={color} {...styles} />
    </svg>
  );
}

export function TallRoundRectangleShape({ color, styles }) {
  return (
    <svg width="41" height="40" viewBox="0 0 32 40" fill="none" xmlns="http://www.w3.org/2000/svg">
      <rect width="32" height="40" rx="5" fill={color} {...styles} />
    </svg>
  );
}

export function BottomRoundRectangleShape({ color, styles }) {
  return (
    <svg width="40" height="40" viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
      <g clipPath="url(#clip0_18568_8265)">
        <g clipPath="url(#clip1_18568_8265)">
          <path
            d="M0 0H40V32.6565C40 36.7177 36.7177 40 32.6565 40H7.34353C3.28234 40 0 36.7177 0 32.6565V0Z"
            fill={color}
            {...styles}
          />
        </g>
      </g>
      <defs>
        <clipPath id="clip0_18568_8265">
          <rect width="40" height="40" fill="none" />
        </clipPath>
        <clipPath id="clip1_18568_8265">
          <rect width="40" height="40" fill="none" />
        </clipPath>
      </defs>
    </svg>
  );
}

export function CutRectangle({ color, styles }) {
  return (
    <svg width="40" height="40" viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
      <g clipPath="url(#clip0_18568_8261)">
        <g clipPath="url(#clip1_18568_8261)">
          <path
            d="M32.4487 0H7.52394L0 7.36623V32.6338L6.97674 40H32.9959L40 32.3836V7.6164L32.4487 0Z"
            fill={color}
            {...styles}
          />
        </g>
      </g>
      <defs>
        <clipPath id="clip0_18568_8261">
          <rect width="40" height="40" fill="none" />
        </clipPath>
        <clipPath id="clip1_18568_8261">
          <rect width="40" height="40" fill="none" />
        </clipPath>
      </defs>
    </svg>
  );
}

export function CircleShape({ color, styles }) {
  return (
    <svg width="40" height="40" viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
      <path
        fillRule="evenodd"
        clipRule="evenodd"
        d="M11.3661 38.0349C1.40555 33.2661 -2.80363 21.3246 1.96473 11.3628C4.25461 6.57899 8.35082 2.90088 13.3523 1.13757C18.3537 -0.62574 23.8507 -0.329761 28.6338 1.96035C38.5946 6.72929 42.8037 18.6709 38.0351 28.6325C33.2669 38.5944 21.3267 42.804 11.3661 38.0349Z"
        fill={color}
        {...styles}
      />
    </svg>
  );
}

export function CylinderShape({ color, styles }) {
  return (
    <svg xmlns="http://www.w3.org/2000/svg" width="40" height="40" viewBox="0 0 34 40" fill="none">
      <path
        d="M1.5 12.0547C2.37103 13.2965 3.7598 14.3745 5.46099 15.2332C8.43989 16.7368 12.5203 17.6523 17 17.6523C21.4797 17.6523 25.5601 16.7368 28.539 15.2332C30.2402 14.3745 31.629 13.2965 32.5 12.0547V30.9238C32.5 32.8651 30.9289 34.7544 28.0884 36.1882C25.2764 37.6075 21.3568 38.5 17 38.5C12.6432 38.5 8.72355 37.6075 5.9116 36.1882C3.07106 34.7544 1.5 32.8651 1.5 30.9238V12.0547Z"
        fill={color}
        stroke="white"
        {...styles}
      />
      <path
        d="M1.5 9.07614C1.5 7.1349 3.07106 5.24559 5.9116 3.81181C8.72355 2.39245 12.6432 1.5 17 1.5C21.3568 1.5 25.2764 2.39245 28.0884 3.81181C30.9289 5.24559 32.5 7.1349 32.5 9.07614C32.5 11.0174 30.9289 12.9067 28.0884 14.3405C25.2764 15.7598 21.3568 16.6523 17 16.6523C12.6432 16.6523 8.72355 15.7598 5.9116 14.3405C3.07106 12.9067 1.5 11.0174 1.5 9.07614Z"
        fill={color}
        stroke="white"
        {...styles}
      />
      <path
        d="M33 9.07614C33 13.5365 25.8366 17.1523 17 17.1523C8.16344 17.1523 1 13.5365 1 9.07614M33 9.07614C33 4.61582 25.8366 1 17 1C8.16344 1 1 4.61582 1 9.07614M33 9.07614V30.9239C33 35.3842 25.8366 39 17 39C8.16344 39 1 35.3842 1 30.9239V9.07614"
        stroke="white"
        {...styles}
      />
    </svg>
  );
}

export function TagShape({ color, styles }) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width="40"
      height="40"
      viewBox="0 0 40 40"
      fill="none"
      x="0px"
      y="0px"
    >
      <g>
        <polygon fill={color} points="0.5,39.3 0.5,0.7 26.4,0.7 39.3,20 26.4,39.3" {...styles} />
        <path
          fill={color}
          d="M26.3,1.3L38.7,20L26.3,38.7H1.1V1.3H26.3 M26.7,0H0v40h26.7L40,20L26.7,0L26.7,0z"
          {...styles}
        />
      </g>
    </svg>
  );
}

export function RoundTagShape({ color, styles }) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width="40"
      height="40"
      viewBox="0 0 40 40"
      fill="none"
      x="0px"
      y="0px"
    >
      <g>
        <path
          fill={color}
          {...styles}
          d="M3,39.3c-1.3,0-2.5-1.5-2.5-3.3V4c0-1.9,1.1-3.3,2.5-3.3h23.2c0.7,0,1.3,0.4,1.7,1.1l11,16
        c0.9,1.3,0.9,3.2,0,4.5l-11,16c-0.5,0.7-1.1,1.1-1.7,1.1H3z"
        />
        <path
          fill={color}
          {...styles}
          d="M26.2,1.3c0.5,0,1.1,0.3,1.5,0.9l11,16c0.7,1.1,0.7,2.5,0,3.6l-11,16c-0.4,0.5-0.9,0.9-1.5,0.9H3
        c-1.1,0-2-1.2-2-2.7V4c0-1.5,0.9-2.7,2-2.7H26.2 M26.2,0H3C1.3,0-0.1,1.7-0.1,4v32c0,2.3,1.3,4,3,4h23.2c0.8,0,1.6-0.5,2.1-1.3
        l11-16c0.9-1.5,0.9-3.9,0-5.3l-10.9-16C27.8,0.5,27,0,26.2,0L26.2,0z"
        />
      </g>
    </svg>
  );
}

export function StarShape({ color, styles }) {
  return (
    <svg width="40" height="40" viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
      <g clipPath="url(#clip0_18568_8307)">
        <path
          d="M19.9869 1L27.4001 12.4366L40 15.5248L31.1722 25.7314L32.351 39L19.9869 34.8388L7.64899 39L8.51343 25.8361L0 15.5248L12.8094 12.5675L19.9869 1Z"
          fill={color}
          {...styles}
        />
      </g>
      <defs>
        <clipPath id="clip0_18568_8307">
          <rect width="40" height="38" fill="none" transform="translate(0 1)" />
        </clipPath>
      </defs>
    </svg>
  );
}

export function VeeShape({ color, styles }) {
  return (
    <svg width="40" height="40" viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
      <g clipPath="url(#clip0_18568_8293)">
        <g clipPath="url(#clip1_18568_8293)">
          <path d="M40 0L20.014 40L0 0L20.014 13.3148L40 0Z" fill={color} {...styles} />
        </g>
      </g>
      <defs>
        <clipPath id="clip0_18568_8293">
          <rect width="40" height="40" fill="none" />
        </clipPath>
        <clipPath id="clip1_18568_8293">
          <rect width="40" height="40" fill="none" />
        </clipPath>
      </defs>
    </svg>
  );
}

export function HexagonShape({ color, styles }) {
  return (
    <svg width="40" height="40" viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
      <g clipPath="url(#clip0_18568_8231)">
        <g clipPath="url(#clip1_18568_8231)">
          <path d="M30 0H10L0 20L10 40H30L40 20L30 0Z" fill={color} {...styles} />
        </g>
      </g>
      <defs>
        <clipPath id="clip0_18568_8231">
          <rect width="40" height="40" fill="none" />
        </clipPath>
        <clipPath id="clip1_18568_8231">
          <rect width="40" height="40" fill="none" />
        </clipPath>
      </defs>
    </svg>
  );
}

export function RoundHexagonShape({ color, styles }) {
  return (
    <svg width="40" height="40" viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
      <g clipPath="url(#clip0_18568_8327)">
        <g clipPath="url(#clip1_18568_8327)">
          <path
            d="M27.3712 0H12.6566C10.8525 0 9.21755 1.00559 8.40007 2.59777L0.535377 17.8492C-0.169344 19.1899 -0.169344 20.8101 0.535377 22.1508L8.40007 37.4022C9.21755 38.9944 10.8525 40 12.6566 40H27.3712C29.1753 40 30.8102 38.9944 31.6277 37.4022L39.4924 22.1508C40.1971 20.8101 40.1971 19.1899 39.4924 17.8492L31.6277 2.59777C30.8102 1.00559 29.1753 0 27.3712 0Z"
            fill={color}
            {...styles}
          />
        </g>
      </g>
      <defs>
        <clipPath id="clip0_18568_8327">
          <rect width="40" height="40" fill="none" />
        </clipPath>
        <clipPath id="clip1_18568_8327">
          <rect width="40" height="40" fill="none" />
        </clipPath>
      </defs>
    </svg>
  );
}

export function ConcaveHexagonShape({ color, styles }) {
  return (
    <svg width="40" height="40" viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
      <g clipPath="url(#clip0_18568_8311)">
        <path d="M35.3963 20L40 39H0L4.63129 20L0 1H40L35.3963 20Z" fill={color} {...styles} />
      </g>
      <defs>
        <clipPath id="clip0_18568_8311">
          <rect width="40" height="38" fill="none" transform="translate(0 1)" />
        </clipPath>
      </defs>
    </svg>
  );
}

export function HeptagonShape({ color, styles }) {
  return (
    <svg width="40" height="40" viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
      <g clipPath="url(#clip0_18568_8331)">
        <g clipPath="url(#clip1_18568_8331)">
          <path
            d="M20 0L3.96648 7.92842L0 25.7398L11.0894 40H28.9106L40 25.7398L36.0335 7.92842L20 0Z"
            fill={color}
            {...styles}
          />
        </g>
      </g>
      <defs>
        <clipPath id="clip0_18568_8331">
          <rect width="40" height="40" fill="none" />
        </clipPath>
        <clipPath id="clip1_18568_8331">
          <rect width="40" height="40" fill="none" />
        </clipPath>
      </defs>
    </svg>
  );
}

export function RoundHeptagonShape({ color, styles }) {
  return (
    <svg width="40" height="40" viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
      <g clipPath="url(#clip0_18568_8319)">
        <g clipPath="url(#clip1_18568_8319)">
          <path
            d="M18.0839 0.445982L5.536 6.63414C4.37205 7.19163 3.54877 8.25086 3.26488 9.50522L0.11371 23.6376C-0.170179 24.8641 0.11371 26.1185 0.9086 27.0941L9.68078 38.3554C10.5041 39.3867 11.7532 40 13.0875 40H26.9412C28.2755 40 29.5246 39.3867 30.3479 38.3554L39.1201 27.0941C39.8866 26.1185 40.1705 24.8362 39.915 23.6376L36.7638 9.47734C36.4799 8.22299 35.6567 7.19163 34.4927 6.60626L21.9164 0.445982C20.6957 -0.139383 19.2762 -0.139383 18.0555 0.445982H18.0839Z"
            fill={color}
            {...styles}
          />
        </g>
      </g>
      <defs>
        <clipPath id="clip0_18568_8319">
          <rect width="40" height="40" fill="none" />
        </clipPath>
        <clipPath id="clip1_18568_8319">
          <rect width="40" height="40" fill="none" />
        </clipPath>
      </defs>
    </svg>
  );
}

export function OctagonShape({ color, styles }) {
  return (
    <svg width="40" height="40" viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
      <g clipPath="url(#clip0_18568_8315)">
        <g clipPath="url(#clip1_18568_8315)">
          <path
            d="M28.2839 0H11.7161L0 11.7161V28.2839L11.7161 40H28.2839L40 28.2839V11.7161L28.2839 0Z"
            fill={color}
            {...styles}
          />
        </g>
      </g>
      <defs>
        <clipPath id="clip0_18568_8315">
          <rect width="40" height="40" fill="none" />
        </clipPath>
        <clipPath id="clip1_18568_8315">
          <rect width="40" height="40" fill="none" />
        </clipPath>
      </defs>
    </svg>
  );
}

export function RoundOctagonShape({ color, styles }) {
  return (
    <svg width="40" height="40" viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
      <g clipPath="url(#clip0_18568_8323)">
        <g clipPath="url(#clip1_18568_8323)">
          <path
            d="M26.6575 0H13.3425C12.295 0 11.3026 0.413508 10.5582 1.15782L1.15782 10.5582C0.413508 11.3026 0 12.295 0 13.3425V26.6575C0 27.705 0.413508 28.6975 1.15782 29.4418L10.5858 38.8697C11.3301 39.6141 12.3225 40.0276 13.3701 40.0276H26.685C27.7326 40.0276 28.725 39.6141 29.4693 38.8697L38.8973 29.4418C39.6416 28.6975 40.0551 27.705 40.0551 26.6575V13.3425C40.0551 12.295 39.6416 11.3026 38.8973 10.5582L29.4418 1.15782C28.6975 0.413508 27.705 0 26.6575 0Z"
            fill={color}
            {...styles}
          />
        </g>
      </g>
      <defs>
        <clipPath id="clip0_18568_8323">
          <rect width="40" height="40" fill="none" />
        </clipPath>
        <clipPath id="clip1_18568_8323">
          <rect width="40" height="40" fill="none" />
        </clipPath>
      </defs>
    </svg>
  );
}

export function PolygonShape({ color, styles }) {
  return (
    <svg width="40" height="40" viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
      <g clipPath="url(#clip0_18568_8289)">
        <g clipPath="url(#clip1_18568_8289)">
          <path
            d="M27.136 13.2153V0H13.1125V13.2153H0V26.8665H13.1125V40H27.136V26.8665H40V13.2153H27.136Z"
            fill={color}
            {...styles}
          />
        </g>
      </g>
      <defs>
        <clipPath id="clip0_18568_8289">
          <rect width="40" height="40" fill="none" />
        </clipPath>
        <clipPath id="clip1_18568_8289">
          <rect width="40" height="40" fill="none" />
        </clipPath>
      </defs>
    </svg>
  );
}

export function RhomboidShape({ color, styles }) {
  return (
    <svg width="40" height="40" viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
      <g clipPath="url(#clip0_18568_8257)">
        <g clipPath="url(#clip1_18568_8257)">
          <path d="M13.2153 40H40L26.812 0H0L13.2153 40Z" fill={color} {...styles} />
        </g>
      </g>
      <defs>
        <clipPath id="clip0_18568_8257">
          <rect width="40" height="40" fill="none" />
        </clipPath>
        <clipPath id="clip1_18568_8257">
          <rect width="40" height="40" fill="none" />
        </clipPath>
      </defs>
    </svg>
  );
}

export function RightRhomboidShape({ color, styles }) {
  return (
    <svg width="40" height="40" viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
      <g clipPath="url(#clip0_18568_8250)">
        <g clipPath="url(#clip1_18568_8250)">
          <path d="M26.812 40H0L13.2153 0H40L26.812 40Z" fill={color} {...styles} />
        </g>
      </g>
      <defs>
        <clipPath id="clip0_18568_8250">
          <rect width="40" height="40" fill="none" />
        </clipPath>
        <clipPath id="clip1_18568_8250">
          <rect width="40" height="40" fill="none" />
        </clipPath>
      </defs>
    </svg>
  );
}

export function DiamondShape({ color, styles }) {
  return (
    <svg width="40" height="40" viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
      <g clipPath="url(#clip0_18568_8243)">
        <g clipPath="url(#clip1_18568_8243)">
          <path
            d="M19.9764 -0.00695638L-0.0166016 19.9861L19.9764 39.9791L39.9695 19.9861L19.9764 -0.00695638Z"
            fill={color}
            {...styles}
          />
        </g>
      </g>
      <defs>
        <clipPath id="clip0_18568_8243">
          <rect width="40" height="40" fill="none" />
        </clipPath>
        <clipPath id="clip1_18568_8243">
          <rect width="40" height="40" fill="none" />
        </clipPath>
      </defs>
    </svg>
  );
}

export function RoundDiamondShape({ color, styles }) {
  return (
    <svg width="40" height="40" viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
      <g clipPath="url(#clip0_18568_8236)">
        <g clipPath="url(#clip1_18568_8236)">
          <path
            d="M16.4412 1.48312L1.49775 16.4266C-0.475803 18.4001 -0.475803 21.5999 1.49775 23.5734L16.4412 38.5169C18.4148 40.4904 21.6145 40.4904 23.5881 38.5169L38.5315 23.5734C40.5051 21.5999 40.5051 18.4001 38.5315 16.4266L23.5881 1.48312C21.6145 -0.490437 18.4148 -0.490435 16.4412 1.48312Z"
            fill={color}
            {...styles}
          />
        </g>
      </g>
      <defs>
        <clipPath id="clip0_18568_8236">
          <rect width="40" height="40" fill="none" />
        </clipPath>
        <clipPath id="clip1_18568_8236">
          <rect width="40" height="40" fill="none" />
        </clipPath>
      </defs>
    </svg>
  );
}

export function BarrelShape({ color, styles }) {
  return (
    <svg width="40" height="40" viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
      <g clipPath="url(#clip0_18628_25765)">
        <path
          d="M39.4419 5.26901V34.7309C39.4419 40.216 30.7255 39.9987 19.9997 39.9987C9.27399 39.9987 0.557617 40.216 0.557617 34.7309V5.26901C0.557617 -0.21606 9.24684 0.00117016 19.9997 0.00117016C30.7527 0.00117016 39.4419 -0.21606 39.4419 5.26901Z"
          fill={color}
          {...styles}
        />
      </g>
      <defs>
        <clipPath id="clip0_18628_25765">
          <rect width="40" height="40" fill="none" />
        </clipPath>
      </defs>
    </svg>
  );
}

export function XWingShape({ color }) {
  return (
    <svg width="40" height="40" viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
      <g>
        <path
          fill={color}
          d="M0,28.7c0.1-0.1,0.2-0.2,0.3-0.3c1.7-1.6,3.3-3.2,5-4.8c1.2-1.1,2.4-2.3,3.5-3.4C9,20,9,20,8.8,19.9
        c-1.2-1.2-2.4-2.3-3.6-3.5c-1.7-1.7-3.5-3.3-5.2-5c0,0,0,0-0.1-0.1c0.1-0.1,0.2-0.2,0.3-0.2c2-1.5,4-2.9,6-4.4
        c0.5-0.3,0.9-0.7,1.4-1c0.1-0.1,0.1,0,0.2,0c1.2,1.4,2.3,2.7,3.5,4.1c1.6,1.9,3.3,3.9,4.9,5.8c1.2,1.4,2.5,2.9,3.7,4.3
        c0.1,0.1,0.1,0.1,0,0.2c-1.8,2.1-3.7,4.3-5.5,6.4c-1.8,2.1-3.7,4.3-5.5,6.4c-0.4,0.4-0.7,0.9-1.1,1.3c-0.1,0.1-0.1,0.1-0.2,0
        c-1.7-1.2-3.3-2.4-5-3.7C1.8,30,1,29.4,0.1,28.8C0.1,28.7,0.1,28.7,0,28.7z"
        />
        <path
          fill={color}
          d="M40,28.7c-0.1,0.1-0.3,0.2-0.4,0.3c-1.4,1-2.8,2-4.2,3c-1,0.7-2.1,1.5-3.1,2.3c-0.1,0.1-0.1,0-0.2,0
        c-1-1.2-2-2.3-3-3.5c-1.4-1.6-2.8-3.2-4.2-4.9c-1.1-1.3-2.2-2.5-3.2-3.8c-0.6-0.7-1.1-1.3-1.7-2C20,20,20,20,20.1,19.9
        c1.4-1.6,2.8-3.3,4.2-4.9c1.2-1.4,2.5-2.9,3.7-4.3c1.2-1.4,2.3-2.7,3.5-4.1C31.7,6.3,32,6,32.2,5.8c0.1-0.1,0.1-0.1,0.2,0
        c1.5,1.1,3,2.2,4.6,3.3c1,0.7,2,1.4,2.9,2.1c0,0,0,0,0.1,0c0.1,0,0.1,0.1,0,0.1c-0.4,0.3-0.7,0.7-1.1,1c-2.4,2.3-4.8,4.6-7.1,6.9
        c-0.2,0.2-0.4,0.4-0.6,0.5C31,20,31,20,31.2,20.1c1.2,1.2,2.4,2.3,3.6,3.5c1.7,1.7,3.4,3.3,5.1,5C39.9,28.6,40,28.6,40,28.7
        C40,28.7,40,28.7,40,28.7z"
        />
      </g>
    </svg>
  );
}
