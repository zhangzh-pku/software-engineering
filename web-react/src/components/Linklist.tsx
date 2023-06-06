import React from 'react';
import { List, ListItem, ListItemText } from '@mui/material';
import Link from '@mui/material/Link';

interface Props {
    links: string[];
}

const LinkList: React.FC<Props> = ({ links }) => {
    return (
        <List>
            {links.map((link) => (
                <ListItem key={link}>
                    {/*<ListItemText primary={link} />*/}
                    <Link href={link} underline='hover' color='primary'>
                        {link}
                    </Link>
                </ListItem>
            ))}
        </List>
    );
};

export default LinkList;