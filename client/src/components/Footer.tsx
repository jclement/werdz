import React, { Component } from 'react';

export class Footer extends Component {
    render() {
        return (
            <div>
                <div className="footer">
                    <div className="footer-links">
                        <div className="footer-attribution">
                            A silly game written by <a href="https://zeos.ca/about/">Jeff Clement</a>
                        </div>
                    </div>
                </div>
            </div>
        );
    }
}