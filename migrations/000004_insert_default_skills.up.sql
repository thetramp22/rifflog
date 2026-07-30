INSERT INTO skills (name, description)
VALUES
    ('Alternate Picking', 'Alternate Picking is a picking technique that alternates downstrokes and upstrokes to improve speed, rhythm, and efficiency.'),
    ('Chord Transition', 'Chord Transition is moving your fretting hand smoothly and efficiently from one chord shape to another in time with the music.'),
    ('Scales', 'A guitar scale is asequence of notes played in ascending or descending order of pitch, separated by specific intervals of whole and half steps.'),
    ('Finger Independence', 'Finger independence is the ability to control and move each finger on your fretting (or picking) hand completely isolated from the others.'),
    ('Ear Training', 'Ear training is the process of learning to recognize and identify musical pitches, intervals, chords, and rhythms by ear.')
ON CONFLICT (name) DO NOTHING;