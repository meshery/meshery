import { MESHERY_CLOUD_PROD } from '@/constants/endpoints';
import Link from '@mui/material/Link';

export const ActionName = {
  NEXT: 'Next',
  SETUP: 'Setup',
  SHARE: 'Share A Design',
  DONE: 'Done',
  EXPLORE: 'Explore',
  CATALOG: 'Choose Template',
  CREATE: 'From Scratch',
  IMPORT: 'Import',
  LEARN: 'Learn',
};

export const stepsData = [
  {
    id: 1,
    isVisit: false,
    title: 'Use the CNCF Playground',
    subTitle:
      "Explore the Cloud Native Computing Foundation's graduated, incubation, and sandbox projects as well as many other popular open source projects.",
    isDisabled: false,
    journey: [
      {
        id: 1,
        title: 'The Cloud Native Playground',
        video: 'https://www.youtube.com/embed/Do7htKrRzDA?si=5iMQ5a1JUf3qpIiH',
        content: (
          <p>
            {
              "Explore the Cloud Native Computing Foundation's graduated, incubation, and sandbox projects as well as many other popular open source projects."
            }
          </p>
        ),
        previousButton: false,
        actionName: ActionName.NEXT,
      },
      {
        id: 2,
        title: 'The Cloud Native Playground',
        image: `/static/img/getting-started/kubernate-clusters.png`,
        content: (
          <>
            <p>
              The Meshery Playground is connected to live Kubernetes cluster(s) and allows users
              full-control over those clusters.
            </p>
            <p>
              <b>Drag-and-drop</b>, your cloud native infrastructure using a palette of thousands of
              versioned Kubernetes components. Say <b>goodbye to YAML configurations.</b>
            </p>
          </>
        ),
        previousButton: true,
        actionName: ActionName.EXPLORE,
        roles: ['meshmap', 'admin'],
        primaryAction: () => {
          window.open('/extension/meshmap?mode=design', '_self');
        },
      },
    ],
  },
  {
    id: 2,
    isVisit: false,
    title: 'GitOps your infra with Kanvas Snapshots',
    subTitle: 'Connect your GitHub repository to receive Kanvas snapshots in your pull requests.',
    isDisabled: false,
    journey: [
      {
        id: 3,
        title: 'Design snapshots in your PRs',
        image: `/static/img/getting-started/kanvas-snapshots.png`,
        content: (
          <p>
            See your deployment before you merge. Connect Kanvas to your GitHub repo and see changes
            pull request-to-pull request.
          </p>
        ),
        previousButton: false,
        actionName: ActionName.SETUP,
        // make it dynmic for cloud and the meshery ui
        primaryAction: () => {
          window.open(`${MESHERY_CLOUD_PROD}/connect/github/new`, '_blank');
        },
      },
    ],
  },
  {
    id: 3,
    isVisit: false,
    title: 'Invite a Friend to Collaborate',
    subTitle: 'Stop finger-pointing and start collaborating.',
    isDisabled: false,
    journey: [
      {
        id: 4,
        title: 'Share your design',
        image: `${process.env.API_ENDPOINT_PREFIX}/assets/images/get-started/share-design.gif`,
        content: (
          <p>
            {
              'Invite your teammates to collaborate on a design. Designs set to public visibility by default, so sharing is easy. If you prefer to share your designs privately be sure to signup for either the Team or Enterprise '
            }
            {<Link href="/account/plans">plan</Link>}
            {'.'}
          </p>
        ),
        previousButton: false,
        actionName: ActionName.SHARE,
        primaryAction: () => {
          window.open('/catalog/content/my-designs', '_self');
        },
      },
    ],
  },
  {
    id: 4,
    isVisit: false,
    title: 'Create a new design',
    subTitle:
      'Either choose a template from the catalog or start from scratch in the Meshery Playground.',
    isDisabled: false,
    journey: [
      {
        id: 4,
        title: 'Create a new design',
        embed: true,
        content: (
          <p>
            To create your first design, save time with a template from the catalog. Alternatively,
            start from scratch with a fresh design.
            <br />
            <p
              style={{
                marginTop: '10px',
                fontWeight: 'bold',
              }}
            >
              Choose your starting point.
            </p>
          </p>
        ),
        previousButton: false,
        actionName: ActionName.NEXT,
      },
      {
        id: 5,
        title: 'Start from scratch',
        image: `/static/img/getting-started/start-from-scratch.gif`,
        content: (
          <p>
            Love a blank canvas? Start by creating your very own design from scratch in the Meshery
            Playground.
          </p>
        ),
        previousButton: true,
        actionName: ActionName.CREATE,
        secondaryActionName: ActionName.NEXT,
        roles: ['meshmap', 'admin'],
        primaryAction: () => {
          window.open('/extension/meshmap?mode=design', '_self');
        },
      },
      {
        id: 6,
        title: 'Choose a design template',
        image: `/static/img/getting-started/select-a-template.gif`,
        content: (
          <p>
            Click on any of the design cards available. It will take you to the information page so
            that you can review the template description and familiarize yourself with its
            components. Ensure that the template&#39;s components meet your specific requirements.
          </p>
        ),
        previousButton: true,
        secondaryActionName: ActionName.NEXT,
        actionName: ActionName.CATALOG,
        primaryAction: () => {
          window.open('/configuration/catalog', '_self');
        },
      },
      {
        id: 7,
        title: 'Import any design',
        image: `/static/img/getting-started/import-design.gif`,
        content: (
          <p>
            {
              'Import your existing designs and existing infrastructure configurations into Meshery. The platform supports a variety of application definition formats.'
            }
          </p>
        ),
        previousButton: true,
        actionName: ActionName.IMPORT,
        primaryAction: () => {
          window.open('/configuration/designs', '_self');
        },
      },
    ],
  },
  {
    id: 5,
    isVisit: false,
    title: 'Mastering Meshery',
    subTitle:
      'Learn all about Meshery from structured learning paths combining theoretical knowledge with hands-on, practical experience.',
    isDisabled: false,
    journey: [
      {
        id: 8,
        title: 'Learning about Meshery',
        image: `/static/img/getting-started/master-meshery-learn.gif`,
        content: (
          <p>
            {
              'Learning Paths are designed to help you understand and master cloud native technologies. Learning Paths offer a structured approach to learning, combining theoretical knowledge with hands-on, practical experience.'
            }
          </p>
        ),
        previousButton: false,
        actionName: ActionName.LEARN,
        primaryAction: () => {
          window.open(`${MESHERY_CLOUD_PROD}/academy/learning-paths/mastering-meshery`, '_blank');
        },
      },
    ],
  },
];
