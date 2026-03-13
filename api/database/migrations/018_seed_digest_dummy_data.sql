-- Migration: Seed digest dummy data
-- Date: 2025-01-15
-- Description: Populates digest tables with dummy data for testing and development

-- Create today's digest
INSERT INTO digests (date, title, summary, is_dummy) VALUES
(CURRENT_DATE, 'Daily News Digest', 'Stay informed with the most important news from around the world.', true);

-- Get today's digest ID
DO $$
DECLARE
    today_digest_id UUID;
BEGIN
    SELECT id INTO today_digest_id FROM digests WHERE date = CURRENT_DATE;

    -- Insert dummy articles for today's digest
    INSERT INTO digest_articles (
        digest_id, title, content, summary, source, author, published_at,
        category, image_url, external_url, read_time_minutes,
        is_breaking, is_hot, is_dummy
    ) VALUES
    (
        today_digest_id,
        'AI Revolution Transforms Healthcare Diagnosis',
        'Revolutionary artificial intelligence algorithms are transforming the way medical professionals diagnose diseases, with new systems achieving 95% accuracy rates in detecting various conditions. These advanced AI models have been trained on millions of medical cases and can now identify patterns that human doctors might miss.

The breakthrough technology combines machine learning with advanced imaging techniques to provide real-time analysis of medical scans, blood tests, and patient symptoms. Leading hospitals worldwide are already implementing these systems, reporting significant improvements in early disease detection.

Dr. Sarah Chen, Chief Medical Officer at the Stanford AI Lab, explains: "This technology doesn''t replace doctors but enhances their capabilities. It''s like having a highly experienced specialist available 24/7 to provide a second opinion."

The AI diagnostic tools have shown particular promise in oncology, cardiology, and neurology. Early trials suggest that cancer detection rates have improved by 40% when AI assistance is used alongside traditional diagnostic methods.

While the technology promises revolutionary changes in healthcare, experts emphasize the importance of maintaining human oversight and ensuring patient privacy in AI-driven medical systems.',
        'New AI algorithms are revolutionizing medical diagnosis with 95% accuracy rates, enhancing doctor capabilities and improving early disease detection.',
        'TechHealth News',
        'Dr. Michael Rodriguez',
        NOW() - INTERVAL '2 hours',
        'Technology',
        'https://example.com/ai-healthcare.jpg',
        'https://techhealth.com/ai-diagnosis-breakthrough',
        4,
        true,
        false,
        true
    ),
    (
        today_digest_id,
        'Climate Summit 2025: Historic Agreement Reached',
        'World leaders at the Climate Summit 2025 have reached a groundbreaking agreement that commits all participating nations to unprecedented climate action. The agreement, signed by representatives from 195 countries, sets binding targets for carbon emission reductions and renewable energy adoption.

The historic pact includes a commitment to achieve net-zero emissions by 2035, a timeline that is five years ahead of previous international agreements. Additionally, developed nations have pledged $500 billion in climate aid to support developing countries in their transition to clean energy.

Key provisions of the agreement include:
• Mandatory 60% reduction in carbon emissions by 2030
• Complete phase-out of coal power by 2032
• Investment of $2 trillion in renewable energy infrastructure
• Protection of 50% of Earth''s land and oceans by 2030

UN Secretary-General Maria Santos called the agreement "a turning point in human history" and emphasized that immediate action is crucial for implementation.

Environmental groups have praised the ambitious targets while noting that success will depend on rigorous monitoring and enforcement mechanisms.',
        'World leaders commit to unprecedented climate action with binding agreements for net-zero emissions by 2035.',
        'Global Climate Network',
        'Elena Petrov',
        NOW() - INTERVAL '4 hours',
        'Environment',
        'https://example.com/climate-summit.jpg',
        'https://climatenews.org/summit-2025-agreement',
        5,
        false,
        true,
        true
    ),
    (
        today_digest_id,
        'Tech Giants Report Record Quarterly Earnings',
        'Major technology companies have announced record-breaking quarterly earnings, with Apple, Google, and Microsoft exceeding analyst expectations by significant margins. The strong performance comes amid continued growth in cloud computing, AI services, and digital transformation initiatives.

Apple reported quarterly revenue of $125 billion, driven by strong iPhone sales and services growth. The company''s Services division, which includes the App Store and Apple Pay, reached an all-time high of $25 billion in revenue.

Google''s parent company Alphabet posted revenues of $95 billion, with Google Cloud growing 35% year-over-year. The company''s AI investments are showing strong returns, particularly in search and advertising technologies.

Microsoft continued its dominance in enterprise software with $78 billion in quarterly revenue. Azure cloud services grew 42%, solidifying Microsoft''s position as a leader in business cloud solutions.

The exceptional performance has led to increased investor confidence in the technology sector, with tech stocks reaching new highs. Analysts predict continued growth as businesses worldwide accelerate their digital transformation efforts.',
        'Apple, Google, and Microsoft exceed expectations with strong quarterly results driven by cloud computing and AI growth.',
        'Financial Times Tech',
        'Robert Kim',
        NOW() - INTERVAL '6 hours',
        'Business',
        'https://example.com/tech-earnings.jpg',
        'https://fintech.com/q4-tech-earnings-record',
        3,
        false,
        false,
        true
    ),
    (
        today_digest_id,
        'SpaceX Achieves New Milestone in Mars Mission',
        'SpaceX has successfully completed a critical test of its Starship vehicle, bringing humanity one step closer to interplanetary travel. The latest mission involved a successful orbital refueling demonstration, a key technology required for long-duration space missions to Mars.

The test mission lasted 48 hours and demonstrated the ability to transfer fuel between two Starship vehicles in orbit. This capability is essential for Mars missions, as the spacecraft will need to refuel in Earth orbit before beginning the journey to the Red Planet.

Elon Musk, SpaceX CEO, announced that the company is now targeting 2027 for the first crewed mission to Mars. "Today''s success brings us significantly closer to making life multiplanetary," Musk stated during a press conference.

NASA Administrator Jennifer Wilson praised the achievement, noting that the technology could also benefit lunar missions and space exploration initiatives. The successful test opens the door for more ambitious space missions and commercial space travel.

The next phase will involve testing the spacecraft''s life support systems and radiation shielding, critical components for ensuring crew safety during the months-long journey to Mars.',
        'SpaceX successfully demonstrates orbital refueling technology, bringing crewed Mars missions closer to reality by 2027.',
        'Space Exploration Daily',
        'Captain Lisa Zhang',
        NOW() - INTERVAL '8 hours',
        'Science',
        'https://example.com/spacex-mars.jpg',
        'https://spaceexploration.com/spacex-mars-milestone',
        4,
        false,
        false,
        true
    ),
    (
        today_digest_id,
        'Renewable Energy Adoption Hits All-Time High',
        'Global renewable energy capacity has reached a historic milestone, with solar and wind power now accounting for 40% of worldwide electricity generation. The International Energy Agency reports that renewable energy installations grew by 30% in the past year, exceeding all previous records.

Solar power led the growth with 180 GW of new capacity added globally, while wind energy contributed 95 GW. The rapid expansion is driven by decreasing costs, improved technology, and strong government support for clean energy initiatives.

China continues to lead in renewable energy deployment, installing more solar capacity than the rest of the world combined. However, significant growth was also seen in India, the United States, and European nations.

The cost of solar electricity has dropped by 85% over the past decade, making it the cheapest source of power in many regions. Wind energy costs have similarly declined, with offshore wind showing particular promise for coastal nations.

Experts predict that renewable energy could account for 60% of global electricity generation by 2030 if current growth trends continue. This rapid transition is crucial for meeting international climate goals and reducing dependence on fossil fuels.',
        'Solar and wind energy now provide 40% of global electricity as renewable capacity reaches record levels.',
        'Energy Today',
        'Dr. Amanda Foster',
        NOW() - INTERVAL '10 hours',
        'Environment',
        'https://example.com/renewable-energy.jpg',
        'https://energytoday.com/renewable-milestone-2025',
        3,
        false,
        false,
        true
    ),
    (
        today_digest_id,
        'Breakthrough Medical Research Promises New Cancer Treatment',
        'Researchers at Johns Hopkins University have announced a groundbreaking discovery in cancer treatment that could revolutionize how we approach the disease. The new therapy, called CAR-T 2.0, has shown remarkable success in clinical trials with a 92% remission rate.

The innovative treatment involves genetically modifying a patient''s own immune cells to better recognize and attack cancer cells. Unlike traditional CAR-T therapy, the new approach can target multiple cancer markers simultaneously, making it effective against a broader range of cancers.

Dr. Rebecca Martinez, lead researcher on the project, explains: "This represents a quantum leap in personalized cancer treatment. We''re essentially training the patient''s immune system to become a more effective cancer-fighting force."

The therapy has shown particular promise in treating blood cancers, brain tumors, and solid organ cancers that were previously considered untreatable. Early results suggest fewer side effects compared to traditional chemotherapy and radiation treatments.

The FDA has granted fast-track designation for the therapy, and researchers hope to begin Phase III trials within six months. If successful, the treatment could be available to patients within two years.',
        'New CAR-T 2.0 therapy shows 92% remission rate in cancer treatment trials, offering hope for previously untreatable cases.',
        'Medical Research Today',
        'Dr. James Patterson',
        NOW() - INTERVAL '12 hours',
        'Health',
        'https://example.com/cancer-research.jpg',
        'https://medresearch.com/cart-2-breakthrough',
        5,
        false,
        false,
        true
    );

END $$;

-- Create a few historical digests for testing pagination
INSERT INTO digests (date, title, summary, is_dummy) VALUES
(CURRENT_DATE - INTERVAL '1 day', 'Yesterday''s News Digest', 'The most important stories from yesterday.', true),
(CURRENT_DATE - INTERVAL '2 days', 'Weekend News Digest', 'Key developments from the weekend.', true),
(CURRENT_DATE - INTERVAL '3 days', 'Weekly News Roundup', 'Major stories from the past week.', true);

-- Add a few articles to historical digests
DO $$
DECLARE
    yesterday_digest_id UUID;
    weekend_digest_id UUID;
BEGIN
    SELECT id INTO yesterday_digest_id FROM digests WHERE date = CURRENT_DATE - INTERVAL '1 day';
    SELECT id INTO weekend_digest_id FROM digests WHERE date = CURRENT_DATE - INTERVAL '2 days';

    -- Yesterday's articles
    INSERT INTO digest_articles (
        digest_id, title, content, summary, source, author, published_at,
        category, read_time_minutes, is_dummy
    ) VALUES
    (
        yesterday_digest_id,
        'Global Markets Show Strong Recovery',
        'Financial markets worldwide have shown remarkable resilience following recent economic uncertainties...',
        'Stock markets rally as investors show confidence in global economic recovery.',
        'Financial News',
        'Sarah Johnson',
        (CURRENT_DATE - INTERVAL '1 day') + TIME '09:00:00',
        'Business',
        3,
        true
    ),
    (
        yesterday_digest_id,
        'New Archaeological Discovery in Egypt',
        'Archaeologists have uncovered a previously unknown tomb in the Valley of the Kings...',
        'Ancient tomb discovered in Egypt provides new insights into pharaonic burial practices.',
        'History Today',
        'Dr. Ahmed Hassan',
        (CURRENT_DATE - INTERVAL '1 day') + TIME '14:30:00',
        'Science',
        4,
        true
    );

    -- Weekend articles
    INSERT INTO digest_articles (
        digest_id, title, content, summary, source, author, published_at,
        category, read_time_minutes, is_dummy
    ) VALUES
    (
        weekend_digest_id,
        'International Sports Championship Results',
        'The weekend saw exciting developments in international sports competitions...',
        'Weekend sports highlights from major international championships.',
        'Sports Weekly',
        'Mike Thompson',
        (CURRENT_DATE - INTERVAL '2 days') + TIME '16:00:00',
        'Sports',
        2,
        true
    );

END $$;